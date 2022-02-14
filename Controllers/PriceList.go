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
)

func isPriceListExisting(self *Models.PriceList) bool {
	collection := DBManager.SystemCollections.PriceList
	filter := bson.M{
		"name": self.Name,
	}
	_, results := Utils.FindByFilter(collection, filter)
	return len(results) > 0
}

func PriceListGetById(id primitive.ObjectID) (Models.PriceList, error) {
	collection := DBManager.SystemCollections.PriceList
	filter := bson.M{"_id": id}
	var self Models.PriceList
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func PriceListGetByIdPopulated(objID primitive.ObjectID, ptr *Models.PriceList) (Models.PriceListPopulated, error) {
	var PriceListDoc Models.PriceList
	if ptr == nil {
		PriceListDoc, _ = PriceListGetById(objID)
	} else {
		PriceListDoc = *ptr
	}
	populatedResult := Models.PriceListPopulated{}
	populatedResult.CloneFrom(PriceListDoc)
	populatedResult.AppliedForContacts = make([]Models.Contact, len(PriceListDoc.AppliedForContacts))
	populatedResult.ProductsList = make([]Models.ProductOfferPopulated, len(PriceListDoc.ProductsList))
	for i, v := range PriceListDoc.AppliedForContacts {
		populatedResult.AppliedForContacts[i], _ = ContactGetById(v)
	}

	for i, v := range PriceListDoc.ProductsList {
		populatedResult.ProductsList[i].CloneFrom(v)
		populatedResult.ProductsList[i].ProductRef, _ = ProductGetById(v.ProductRef)
	}

	return populatedResult, nil
}
func PriceListSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.PriceList
	if c.Params("id") == "" || c.Params("new_status") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	newValue := true
	if c.Params("new_status") == "inactive" {
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
		return errors.New("an error occurred when modifing PriceList status")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func PriceListModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.PriceList
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	var self Models.PriceList
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing PriceListdocument")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func pricelistGetAll(self *Models.PriceListSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.PriceList
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetPriceListSearchBSONObj())
	if !b {
		return results, errors.New("no object found")
	}
	return results, nil
}

func PriceListGetAll(c *fiber.Ctx) error {
	var self Models.PriceListSearch
	c.BodyParser(&self)
	results, err := pricelistGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func PriceListGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.PriceList
	var self Models.PriceListSearch
	c.BodyParser(&self)
	b, results := Utils.FindByFilter(collection, self.GetPriceListSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}

	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.PriceList
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]Models.PriceListPopulated, len(ResultDocs))
	for i, v := range ResultDocs {
		populatedResult[i], _ = PriceListGetByIdPopulated(v.ID, &v)
	}

	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func PriceListNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.PriceList
	var self Models.PriceList
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	if isPriceListExisting(&self) {
		return errors.New("obj is already Found")
	}
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}
