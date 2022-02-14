package Controllers

import (
	"context"
	"encoding/json"
	"errors"

	"SEEN-TECH-VAI21-BACKEND-GO/DBManager"
	"SEEN-TECH-VAI21-BACKEND-GO/Models"
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UnitsOfMeasurementNew(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.UnitsOfMeasurement

	// Fill the received data inside an obj
	var self Models.UnitOfMeasurement
	c.BodyParser(&self)

	// Validate the obj
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	_, err = collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}

	c.Set("Content-Type", "application/json")
	c.Status(200).Send([]byte("Added Successfully"))

	return nil
}

func UnitsOfMeasurementAddRelations(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.UnitsOfMeasurement

	// Fill the received data inside an obj
	var self []Models.UOMRelation
	c.BodyParser(&self)

	if c.Params("to_id") == "" {
		c.Status(404)
		return errors.New("invalid request params")
	}
	targetID, _ := primitive.ObjectIDFromHex(c.Params("to_id"))

	// https://docs.mongodb.com/manual/reference/operator/update/push/#mongodb-update-up.-push
	updateData := bson.M{
		"$push": Models.UOMRelationGetAppendBSONObj(self),
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": targetID}, updateData)

	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing storage condition status")
	} else {
		c.Status(200)
	}

	return nil
}

func UnitsOfMeasurementGetAll(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.UnitsOfMeasurement

	// Fill the received search obj data
	var self Models.UnitOfMeasurementSearch
	c.QueryParser(&self)

	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetUnitOfMeasurementSearchBSONObj())
	if !b {
		err := errors.New("db error")
		c.Status(500).Send([]byte(err.Error()))
		return err
	}

	// Decode
	response, _ := json.Marshal(
		bson.M{"result": results},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)

	return nil
}

func UnitsOfMeasurementGetDistinctCategories(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.UnitsOfMeasurement

	results, err := collection.Distinct(context.Background(), "category", bson.M{})
	if err != nil {
		c.Status(500)
		return err
	}

	// Decode
	response, _ := json.Marshal(
		bson.M{"result": results},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)

	return nil
}

func UnitsOfMeasurementGetById(objID primitive.ObjectID) (Models.UnitOfMeasurement, error) {

	var self Models.UnitOfMeasurement

	var filter bson.M = bson.M{}
	filter = bson.M{"_id": objID}

	collection := DBManager.SystemCollections.UnitsOfMeasurement

	var results []bson.M
	b, results := Utils.FindByFilter(collection, filter)
	if !b || len(results) == 0 {
		return self, errors.New("obj not found")
	}

	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode

	return self, nil
}

func UnitsOfMeasurementGetAllPopulated(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.UnitsOfMeasurement

	// Fill the received search obj data
	var self Models.UnitOfMeasurementSearch
	c.QueryParser(&self)

	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetUnitOfMeasurementSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("db err or obj is not exist")
	}

	// Decode
	UoMsResult, _ := json.Marshal(results)

	var UoMDocs []Models.UnitOfMeasurement
	json.Unmarshal(UoMsResult, &UoMDocs) //Encode

	populatedResult := make([]Models.UnitOfMeasurementPopulated, len(UoMDocs))

	// Populated the storage Condition
	for i, v := range UoMDocs {
		populatedResult[i].CloneFrom(v)

		for _, w := range v.Relations {
			embeddedUoMDoc, _ := UnitsOfMeasurementGetById(w.UnitRef)
			var populatedEmbeddedUomDoc Models.UOMRelationPopulated
			populatedEmbeddedUomDoc.Ratio = w.Ratio
			populatedEmbeddedUomDoc.Status = w.Status
			populatedEmbeddedUomDoc.UnitRef = embeddedUoMDoc
			populatedResult[i].Relations = append(populatedResult[i].Relations, populatedEmbeddedUomDoc)
		}
	}

	allpopulated, _ := json.Marshal(populatedResult)
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)

	return nil
}

func UnitsOfMeasurementGetDistinctUOM(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.UnitsOfMeasurement

	results, err := collection.Distinct(context.Background(), "name", bson.M{})
	if err != nil {
		c.Status(500)
		return err
	}

	// Decode
	response, _ := json.Marshal(
		bson.M{"result": results},
	)
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)

	return nil
}

func UnitsOfMeasurementSetStatus(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.UnitsOfMeasurement

	if c.Params("id") == "" || c.Params("new_status") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	var newValue = true
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
		return errors.New("an error occurred when modifing unit of measurements status")
	}

	c.Status(200).Send([]byte("status modified successfully"))
	return nil
}

func UnitsOfMeasurementSetRelationStatus(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.UnitsOfMeasurement

	if c.Params("id") == "" || c.Params("embed_id") == "" || c.Params("new_status") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	embedID, _ := primitive.ObjectIDFromHex(c.Params("embed_id"))

	var newValue = true
	if c.Params("new_status") == "inactive" {
		newValue = false
	}

	updateData := bson.M{
		"$set": bson.M{
			"relations.$.status": newValue,
		},
	}

	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID, "relations.unitref": embedID}, updateData)

	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifing storage condition status")
	}

	c.Status(200).Send([]byte("relation status modified successfully"))
	return nil
}

func unitsOfMeasurementConvert(fromAmount float64, fromRef primitive.ObjectID, toRef primitive.ObjectID) (float64, error) {
	var self Models.UnitOfMeasurement
	filter := bson.M{
		"_id":               fromRef,
		"relations.unitref": toRef,
	}

	collection := DBManager.SystemCollections.UnitsOfMeasurement

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)

	if len(results) == 0 {
		// try Backward
		filter := bson.M{
			"_id":               toRef,
			"relations.unitref": fromRef,
		}
		var resultsBackward []bson.M
		cur, err := collection.Find(context.Background(), filter)
		if err != nil {
			return 0, err
		}
		defer cur.Close(context.Background())

		cur.All(context.Background(), &resultsBackward)
		if len(resultsBackward) == 0 {
			return 0, errors.New("cann't convert. relation not found")
		}

		bsonBytes, _ := bson.Marshal(resultsBackward[0]) // Decode
		bson.Unmarshal(bsonBytes, &self)
		for _, v := range self.Relations {
			if v.UnitRef == fromRef {
				return 1 / v.Ratio * fromAmount, nil
			}
		}
		return 0, errors.New("couldn't convert to target UoM")
	}

	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode

	for _, v := range self.Relations {
		if v.UnitRef == toRef {
			return v.Ratio * fromAmount, nil
		}
	}

	return 0, errors.New("couldn't convert to target UoM")
}

func UnitsOfMeasurementConvertEP(c *fiber.Ctx) error {
	var self Models.UOMConvert
	c.BodyParser(&self)

	convertedAmount, err := unitsOfMeasurementConvert(self.Amount, self.FromUnitRef, self.ToUnitRef)
	if err != nil {
		c.Status(500)
		return err
	}

	response, _ := json.Marshal(
		bson.M{"result": convertedAmount},
	)

	c.Status(200).Send(response)
	return nil
}
