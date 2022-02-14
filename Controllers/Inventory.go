package Controllers

import (
	"SEEN-TECH-VAI21-BACKEND-GO/DBManager"
	"SEEN-TECH-VAI21-BACKEND-GO/Models"
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"
	"context"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func isInventoryIdExisting(objId primitive.ObjectID) bool {
	collection := DBManager.SystemCollections.Inventory
	filter := bson.M{"_id": objId}
	var results []bson.M
	b, results := Utils.FindByFilter(collection, filter)
	return (b && len(results) > 0)
}

func isInventoryNameExisting(collection *mongo.Collection, InventoryName string) bool {

	var filter bson.M = bson.M{
		"inventoryname": InventoryName,
	}
	var results []bson.M
	b, results := Utils.FindByFilter(collection, filter)
	return (b && len(results) > 0)
}

func InventoryNew(c *fiber.Ctx) error {
	Collection := DBManager.SystemCollections.Inventory

	var self Models.Inventory
	c.BodyParser(&self)

	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	if isInventoryNameExisting(Collection, self.InventoryName) {
		c.Status(500)
		return errors.New("inventory Name is already exist")
	}

	//added for make an array of Contents and check here for fronted value
	if len(self.Contents) == 0 {
		self.Contents = []Models.InventoryContentRow{}
	}

	res, err := Collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return errors.New("an error occurred when adding new inventory")

	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)

	return nil

}

func InventoryGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Inventory

	results := []bson.M{}

	var searchParams Models.InventorySearch
	c.BodyParser(&searchParams)

	cur, err := collection.Find(context.Background(), searchParams.GetBSONSearchObj())
	if err != nil {
		c.Status(500)
		return err
	}
	defer cur.Close(context.Background())
	cur.All(context.Background(), &results)
	response, _ := json.Marshal(bson.M{
		"results": results,
	})
	c.Set("content-type", "application/json")
	c.Status(200).Send(response)

	return nil
}
func InventorySetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Inventory
	if c.Params("id") == "" || c.Params("new_status") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	var newValue = true
	if c.Params("new_status") == "Inactive" {
		newValue = false
	}
	updateData := bson.M{
		"$set": bson.M{
			"status": newValue,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing inventory status")
	}
	c.Status(200).Send([]byte("status modified successfully"))
	return nil
}

func InventoryModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Inventory
	var self Models.Inventory
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	updateQuery := bson.M{
		"$set": self.GetBSONModificationObj(),
	}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": self.ID}, updateQuery)
	if err != nil {
		c.Status(500)
		return err
	} else {
		c.Status(200)
	}
	return nil
}

func InventoryGetById(id primitive.ObjectID) (Models.Inventory, error) {
	collection := DBManager.SystemCollections.Inventory
	filter := bson.M{"_id": id}
	var self Models.Inventory
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func InventoryGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Inventory
	var self Models.InventorySearch
	c.BodyParser(&self)
	b, results := Utils.FindByFilter(collection, self.GetBSONSearchObj())
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Inventory
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]Models.InventoryPopulated, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i], _ = InventoryGetByIdPopulated(v.ID, &v)
	}
	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func InventoryGetByIdPopulated(objID primitive.ObjectID, ptr *Models.Inventory) (Models.InventoryPopulated, error) {
	var InventoryDoc Models.Inventory
	if ptr == nil {
		InventoryDoc, _ = InventoryGetById(objID)
	} else {
		InventoryDoc = *ptr
	}

	populatedResult := Models.InventoryPopulated{}
	populatedResult.CloneFrom(InventoryDoc)

	var err1 error
	var err2 error

	// populate for Inventory of Contents array
	populatedResult.Contents = make([]Models.InventoryContentRowPopulated, len(InventoryDoc.Contents))
	for i := range InventoryDoc.Contents {
		populatedResult.Contents[i].CloneFrom(InventoryDoc.Contents[i])
		populatedResult.Contents[i].RawMaterialRef, err1 = MaterialGetById(InventoryDoc.Contents[i].RawMaterialRef)
		populatedResult.Contents[i].ProductRef, err2 = ProductGetById(InventoryDoc.Contents[i].ProductRef)
		if err1 != nil && err2 != nil {
			return populatedResult, err1
		}
	}

	return populatedResult, nil
}

// for add new Content in Inventory
func InventorySetOnHand(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Inventory
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id is not sent")
	}
	inventoryId, _ := primitive.ObjectIDFromHex(c.Params("id")) // get inventory id
	InventoryObject, _ := InventoryGetById(inventoryId)

	var newContent []Models.InventoryContentRow
	c.BodyParser(&newContent)

	// need to validate here for content row using loop
	for _, row := range newContent {
		err := row.Validate()
		if err != nil {
			return err
		}
	}

	arr := make([]primitive.ObjectID, 0)

	// check every content in inventory
	if InventoryObject.Type == "Product Inventory" {
		for _, row := range InventoryObject.Contents { // make an array of primitive.ObjectID
			arr = append(arr, row.ProductRef)
		}
		for _, singleContent := range newContent {
			checkRow := Utils.Contains(arr, singleContent.ProductRef)
			if !checkRow {
				InventoryObject.Contents = append(InventoryObject.Contents, singleContent)
			} else {
				return errors.New("this product already exists")
			}
		}
	} else if InventoryObject.Type == "Raw Material Inventory" {
		for _, row := range InventoryObject.Contents { // make an array of primitive.ObjectID
			arr = append(arr, row.RawMaterialRef)
		}
		for _, singleContent := range newContent {
			checkRow := Utils.Contains(arr, singleContent.RawMaterialRef)
			if !checkRow {
				InventoryObject.Contents = append(InventoryObject.Contents, singleContent)
			} else {
				return errors.New("this raw material already exist")
			}
		}
	}

	newValue := InventoryObject.Contents
	filter := bson.M{"_id": inventoryId}
	updateData := bson.M{
		"$set": bson.M{
			"contents": newValue,
		},
	}

	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when adding new Content")
	} else {
		c.Status(200).Send([]byte("Added New Content In the Inventory Successfully"))
	}

	return nil
}
