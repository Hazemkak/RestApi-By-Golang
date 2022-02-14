package Controllers

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Models"
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"
	"context"
	"encoding/json"
	"errors"

	"SEEN-TECH-VAI21-BACKEND-GO/DBManager"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func isMaterialNameExisting(name string) (bool, interface{}) {
	collection := DBManager.SystemCollections.Material
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

func MaterialNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Material
	var self Models.Material
	c.BodyParser(&self)

	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	_, existing := isMaterialNameExisting(self.Name)
	if existing != "" {
		return errors.New("material name already exists with same name")
	}

	res, err := collection.InsertOne(context.Background(), self)

	if err != nil {
		c.Status(500)
		return err
	}
	responese, _ := json.Marshal(res)
	c.Status(200).Send(responese)

	return nil
}

func MaterialGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Material
	results := []bson.M{}
	cur, err := collection.Find(context.Background(), bson.M{})
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

func MaterialModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Material
	var self Models.Material
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

func MaterialSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Material

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

func MaterialGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Material

	b, results := Utils.FindByFilter(collection, bson.M{})
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	byteArr, _ := json.Marshal(results)
	var Results []Models.Material
	json.Unmarshal(byteArr, &Results)
	populatedResult := make([]Models.MaterialPopulated, len(Results))
	for i, v := range Results {
		populatedResult[i], _ = MaterialGetByIdPopulated(v.ID, &v)
	}
	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func MaterialGetByIdPopulated(objID primitive.ObjectID, ptr *Models.Material) (Models.MaterialPopulated, error) {
	var ProductDoc Models.Material
	if ptr == nil {
		ProductDoc, _ = MaterialGetById(objID)
	} else {
		ProductDoc = *ptr
	}
	populatedResult := Models.MaterialPopulated{}
	populatedResult.CloneFrom(ProductDoc)
	var err error

	// populate for diagramuomid
	if ProductDoc.DiameterUomId != primitive.NilObjectID {
		populatedResult.DiameterUomId, err = UnitsOfMeasurementGetById(ProductDoc.DiameterUomId)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate for weightuomidt
	if ProductDoc.WeightUomId != primitive.NilObjectID {
		populatedResult.WeightUomId, err = UnitsOfMeasurementGetById(ProductDoc.WeightUomId)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate for LengthUomId
	if ProductDoc.LengthUomId != primitive.NilObjectID {
		populatedResult.LengthUomId, err = UnitsOfMeasurementGetById(ProductDoc.LengthUomId)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate for ThicknessUomId
	if ProductDoc.ThicknessUomId != primitive.NilObjectID {
		populatedResult.ThicknessUomId, err = UnitsOfMeasurementGetById(ProductDoc.ThicknessUomId)
		if err != nil {
			return populatedResult, err
		}
	}

	return populatedResult, nil
}

func MaterialGetById(id primitive.ObjectID) (Models.Material, error) {
	collection := DBManager.SystemCollections.Material
	filter := bson.M{"_id": id}
	var self Models.Material
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}
