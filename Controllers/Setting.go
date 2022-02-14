package Controllers

import (
	"SEEN-TECH-VAI21-BACKEND-GO/DBManager"
	"SEEN-TECH-VAI21-BACKEND-GO/Models"
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InitializeSetting() bool {
	collection := DBManager.SystemCollections.Setting
	_, results := Utils.FindByFilter(collection, bson.M{})
	if len(results) <= 0 { //no settings is initialized
		var self Models.Setting
		self.ID = primitive.NewObjectID()
		self.ProductSerial = 1
		_, err := collection.InsertOne(context.Background(), self)
		if err != nil {
			return false
		}
		fmt.Println("Initializing Setting Is Done")
	}
	return true
}

func SettingNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Setting
	var self Models.Setting
	c.BodyParser(&self)
	res, err := collection.InsertOne(context.Background(), self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(res)
	c.Status(200).Send(response)
	return nil
}

func settingGetAll(self *Models.SettingSearch) ([]bson.M, error) {
	collection := DBManager.SystemCollections.Setting
	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetSettingSearchBSONObj())
	if !b || len(results) <= 0 {
		return results, errors.New("no settings object found")
	}
	return results, nil
}

func SettingGetAll(c *fiber.Ctx) error {
	var self Models.SettingSearch
	c.BodyParser(&self)
	results, err := settingGetAll(&self)
	if err != nil {
		c.Status(500)
		return err
	}
	response, _ := json.Marshal(bson.M{"result": results})
	c.Set("Content-Type", "application/json")
	c.Status(200).Send(response)
	return nil
}
