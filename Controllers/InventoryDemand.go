package Controllers

import (
	"context"
	"encoding/json"
	"errors"

	"time"

	"SEEN-TECH-VAI21-BACKEND-GO/DBManager"
	"SEEN-TECH-VAI21-BACKEND-GO/Models"
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InventoryDemandGetById(id primitive.ObjectID) (Models.InventoryDemand, error) {
	collection := DBManager.SystemCollections.InventoryDemand
	filter := bson.M{"_id": id}
	var self Models.InventoryDemand
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) <= 0 {
		return self, errors.New("obj not found")
	}
	if len(results) > 1 {
		return self, errors.New("obj id isn't unique")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func InventoryDemandNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.InventoryDemand

	var self Models.InventoryDemand
	c.BodyParser(&self)

	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	if !isInventoryIdExisting(self.FromInventory) || !isInventoryIdExisting(self.ToInventory) {
		return errors.New("to or from inventory cannot be undefined")
	}

	for _, demand := range self.DemandList {
		if demand.Amount <= 0 {
			return errors.New("product/material amount cannot be zero or below")
		}
		if !demand.ProductRef.IsZero() {
			_, err = ProductGetById(demand.ProductRef)
		} else {
			_, err = MaterialGetById(demand.MaterialRef)
		}

		if err != nil {
			return errors.New("a product/material or more isn't defined")
		}
	}

	self.DateTime = primitive.NewDateTimeFromTime(time.Now())
	self.RequestStatus = true

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}

func InventoryDemandGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.InventoryDemand

	var results []bson.M

	var searchRequests Models.InventoryDemandSearch
	c.BodyParser(&searchRequests)

	b, results := Utils.FindByFilter(collection, searchRequests.GetInventoryDemandSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("no obj found")
	}
	response, _ := json.Marshal(bson.M{"results": results})
	c.Set("content-type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func InventoryDemandGetPopulatedById(objID primitive.ObjectID, ptr *Models.InventoryDemand) (Models.InventoryDemandPopulated, error) {
	var currentDoc Models.InventoryDemand
	if ptr == nil {
		currentDoc, _ = InventoryDemandGetById(objID)
	} else {
		currentDoc = *ptr
	}

	populatedResult := Models.InventoryDemandPopulated{}
	populatedResult.CloneFrom(currentDoc)
	populatedResult.FromInventory, _ = InventoryGetById(currentDoc.FromInventory)
	populatedResult.ToInventory, _ = InventoryGetById(currentDoc.ToInventory)
	//TODO: add populated userRef
	allProductPopulated := make([]Models.DemandContentPopulated, len(currentDoc.DemandList))
	for i, content := range currentDoc.DemandList {
		if !content.ProductRef.IsZero() {
			allProductPopulated[i].ProductRef, _ = ProductGetById(content.ProductRef)
		} else {
			allProductPopulated[i].MaterialRef, _ = MaterialGetById(content.MaterialRef)
		}
		allProductPopulated[i].Amount = content.Amount
		allProductPopulated[i].Sent = content.Sent
		allProductPopulated[i].Recieved = content.Recieved
	}
	populatedResult.DemandList = allProductPopulated
	populatedResult.Type = currentDoc.Type
	populatedResult.RequestStatus = currentDoc.RequestStatus
	return populatedResult, nil
}

func InventoryDemandGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.InventoryDemand

	var results []bson.M
	var searchRequests Models.InventoryDemandSearch
	c.BodyParser(&searchRequests)

	b, results := Utils.FindByFilter(collection, searchRequests.GetInventoryDemandSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("no obj found")
	}

	byteArr, err := json.Marshal(results)
	if err != nil {
		return errors.New("error while getting results")
	}
	var allRequestsDocuments []Models.InventoryDemand
	json.Unmarshal(byteArr, &allRequestsDocuments)

	populatetedResults := make([]Models.InventoryDemandPopulated, len(allRequestsDocuments))

	for i, v := range allRequestsDocuments {
		populatetedResults[i], _ = InventoryDemandGetPopulatedById(v.ID, &v)
	}

	allpopulated, _ := json.Marshal(bson.M{"results": populatetedResults})

	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func InventoryDemandReject(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return errors.New("params are not sent correctly")
	}

	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	collection := DBManager.SystemCollections.InventoryDemand

	_, err := InventoryDemandGetById(objId)
	if err != nil {
		return err
	}

	updateData := bson.M{
		"$set": bson.M{
			"requeststatus": false,
			"enddate":       primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objId}, updateData)
	if err != nil {
		return errors.New("error while rejecting inventory request")
	}
	c.Status(200).Send([]byte("rejected successfully"))
	return nil
}

func InventoryDemandDelete(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return errors.New("params are not sent correctly")
	}

	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	collection := DBManager.SystemCollections.InventoryDemand
	inventoryDemand, err := InventoryDemandGetById(objId)
	if err != nil {
		return err
	}

	if !inventoryDemand.RequestStatus || inventoryDemand.Type != "Request" {
		return errors.New("can't delete request while its rejected or proceeded in the transfer and delivery process")
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objId})
	if err != nil {
		return errors.New("error while rejecting inventory request")
	}
	c.Status(200).Send([]byte("deleted successfully"))
	return nil
}

// Purpose: transfer from Request State to Order State
func InventoryDemandOrderPhase(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return errors.New("params are not sent correctly")
	}

	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	collection := DBManager.SystemCollections.InventoryDemand

	currentDoc, err := InventoryDemandGetPopulatedById(objId, nil)
	if err != nil {
		return err
	}

	if !currentDoc.RequestStatus {
		return errors.New("cannot modify a rejected inventory demand")
	}

	// Check Sent Amounts
	var demandContentArr []Models.DemandContent
	c.BodyParser(&demandContentArr)

	for _, demand := range demandContentArr {
		// check if sent qty is negative
		if demand.Sent < 0 {
			return errors.New("sent amount can't be negative number")
		}
		for _, content := range currentDoc.ToInventory.Contents {
			if currentDoc.ToInventory.Type == "Raw Material Inventory" {
				if demand.MaterialRef == content.RawMaterialRef && demand.Sent > content.Amount {
					return errors.New("can't send more than amount available in stock")
				}
			} else {
				if demand.ProductRef == content.ProductRef && demand.Sent > content.Amount {
					return errors.New("can't send more than amount available in stock")
				}
			}
		}
	}

	updateData := bson.M{
		"$set": bson.M{
			"demandlist":    demandContentArr,
			"type":          "Order",
			"frominventory": currentDoc.ToInventory.ID,
			"toinventory":   currentDoc.FromInventory.ID,
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objId}, updateData)
	if err != nil {
		return errors.New("error while updating inventory request")
	}

	c.Status(200).Send([]byte("modified successfully"))
	return nil
}

// Purpose: transfer from Order State to Delivery State
func InventoryDemandDeliveryPhase(c *fiber.Ctx) error {
	if c.Params("id") == "" {
		return errors.New("params are not sent correctly")
	}

	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	collection := DBManager.SystemCollections.InventoryDemand
	inventoryDemand, err := InventoryDemandGetPopulatedById(objId, nil)
	if err != nil {
		return err
	}

	if inventoryDemand.Type != "Order" {
		return errors.New("request must be in order transfer process")
	}

	// Check Recieved Amounts
	var demandContentArr []Models.DemandContent
	c.BodyParser(&demandContentArr)
	for _, demand := range demandContentArr {
		// check if recieved qty is negative
		if demand.Recieved < 0 {
			return errors.New("recieved amount can't be negative number")
		}
		if demand.Recieved > demand.Sent {
			return errors.New("recieved amount can't be greater than sent amount")
		}
	}

	updateData := bson.M{
		"$set": bson.M{
			"demandlist": demandContentArr,
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objId}, updateData)
	if err != nil {
		return errors.New("error while updating inventory request recieved amount")
	}

	found := true
	// Decrement supplier
	for _, demand := range demandContentArr {
		if !found {
			return errors.New("one product/material or more doen't exist in supplier inventory")
		}
		found = false
		for j, content := range inventoryDemand.FromInventory.Contents {
			if inventoryDemand.FromInventory.Type == "Raw Material Inventory" {
				if demand.MaterialRef == content.RawMaterialRef {
					inventoryDemand.FromInventory.Contents[j].Amount -= demand.Recieved
					found = true
					break
				}
			} else {

				if demand.ProductRef == content.ProductRef {
					inventoryDemand.FromInventory.Contents[j].Amount -= demand.Recieved
					found = true
					break
				}
			}
		}
	}

	filter := bson.M{
		"_id": inventoryDemand.FromInventory.ID,
	}

	updateData = bson.M{
		"$set": bson.M{
			"contents": inventoryDemand.FromInventory.Contents,
		},
	}
	_, err = DBManager.SystemCollections.Inventory.UpdateOne(context.Background(), filter, updateData)
	if err != nil {
		return errors.New("error while decrementing amount from supplier")
	}

	// Increment demander
	for _, demand := range demandContentArr {
		found := false
		for j, content := range inventoryDemand.ToInventory.Contents {
			if inventoryDemand.ToInventory.Type == "Raw Material Inventory" {
				if demand.MaterialRef == content.RawMaterialRef {
					found = true
					inventoryDemand.ToInventory.Contents[j].Amount += demand.Recieved
					break
				}
			} else {
				if demand.ProductRef == content.ProductRef {
					found = true
					inventoryDemand.ToInventory.Contents[j].Amount += demand.Recieved
					break
				}
			}
		}
		if !found {
			inventoryDemand.ToInventory.Contents = append(inventoryDemand.ToInventory.Contents, Models.InventoryContentRow{Amount: demand.Recieved, ProductRef: demand.ProductRef, RawMaterialRef: demand.MaterialRef})
		}
	}

	filter = bson.M{
		"_id": inventoryDemand.ToInventory.ID,
	}

	updateData = bson.M{
		"$set": bson.M{
			"contents": inventoryDemand.ToInventory.Contents,
		},
	}

	_, err = DBManager.SystemCollections.Inventory.UpdateOne(context.Background(), filter, updateData)
	if err != nil {
		return errors.New("error while incrementing amount to demander")
	}

	updateData = bson.M{
		"$set": bson.M{
			"type":    "Delivery",
			"enddate": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": objId}, updateData)
	if err != nil {
		return errors.New("error while Delivering inventory request")
	}

	c.Status(200).Send([]byte("delivered successfully"))
	return nil
}
