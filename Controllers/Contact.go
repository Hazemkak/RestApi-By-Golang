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

func isContactExisting(self *Models.Contact) bool {
	collection := DBManager.SystemCollections.Contact
	filter := bson.M{
		"name": self.Name,
	}
	_, results := Utils.FindByFilter(collection, filter)
	return len(results) > 0
}

func ContactGetById(id primitive.ObjectID) (Models.Contact, error) {
	collection := DBManager.SystemCollections.Contact
	filter := bson.M{"_id": id}
	var self Models.Contact
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

func ContactGetByIdPopulated(objID primitive.ObjectID, ptr *Models.Contact) (Models.ContactPopulated, error) {
	var ContactDoc Models.Contact
	if ptr == nil {
		ContactDoc, _ = ContactGetById(objID)
	} else {
		ContactDoc = *ptr
	}
	populatedResult := Models.ContactPopulated{}
	populatedResult.CloneFrom(ContactDoc)
	populatedResult.CompanyRef, _ = ContactGetById(ContactDoc.CompanyRef)
	return populatedResult, nil
}

func ContactSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Contact
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
		return errors.New("an error occurred when modifing Contact status")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func ContactModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Contact
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{
		"_id": objID,
	}
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		c.Status(404)
		return errors.New("id is not found")
	}
	var self Models.Contact
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
		return errors.New("an error occurred when modifing Contactdocument")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func contactGetAll(self *Models.ContactSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Contact
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetContactSearchBSONObj())
	if !b {
		return results, errors.New("No object found")
	}
	return results, nil
}

func ContactGetAll(c *fiber.Ctx) error {
	var self Models.ContactSearch
	c.BodyParser(&self)
	results, err := contactGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}

func ContactGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Contact
	var self Models.ContactSearch
	c.BodyParser(&self)
	b, results := Utils.FindByFilter(collection, self.GetContactSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Contact
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]Models.ContactPopulated, len(ResultDocs))
	for i, v := range ResultDocs {
		populatedResult[i], _ = ContactGetByIdPopulated(v.ID, &v)
	}
	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func ContactNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Contact
	var self Models.Contact
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	if isContactExisting(&self) {
		return errors.New("Obj is already Found")
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
