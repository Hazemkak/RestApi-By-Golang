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

func ProdStagesNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ProdStages

	var self Models.ProdStages
	c.BodyParser(&self)

	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	_, existing := isProdStagesNameExisting(self.Name)
	if existing != "" {
		return errors.New("prod name already exists with same name")
	}

	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		return err
	}

	response, _ := json.Marshal(res)
	c.Status(200).Send(response)

	return nil
}

func isProdStagesNameExisting(name string) (bool, interface{}) {
	collection := DBManager.SystemCollections.ProdStages
	filter := bson.M{
		"name": name,
	}

	b, results := Utils.FindByFilter(collection, filter)
	id := ""
	if len(results) > 0 {
		id = results[0]["_id"].(primitive.ObjectID).Hex()
	}
	return b, id
}

func ProdStagesGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ProdStages
	results := []bson.M{}

	var searchParams Models.ProdStageSearch
	c.BodyParser(&searchParams)

	cur, err := collection.Find(context.Background(), searchParams.GetProdStagesSearchBSONObj())
	if err != nil {
		c.Status(500)
		return err
	}
	defer cur.Close(context.Background())

	cur.All(context.Background(), &results)
	response, _ := json.Marshal(bson.M{
		"result": results,
	})

	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func ProdStagesModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ProdStages
	var self Models.ProdStages
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

func ProdStagesSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.ProdStages

	if c.Params("id") == "" || c.Params("new_status") == "" {
		c.Status(404)
		return errors.New("all params not send correctly")
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
		return errors.New("an error occurred when mofifing material status")
	}

	c.Status(200).Send([]byte("Status Modified Successfully"))
	return nil
}
