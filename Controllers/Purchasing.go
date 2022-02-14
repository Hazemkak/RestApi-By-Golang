package Controllers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"SEEN-TECH-VAI21-BACKEND-GO/DBManager"
	"SEEN-TECH-VAI21-BACKEND-GO/Models"
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PurchasingNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Purchasing
	var self Models.Purchasing
	c.BodyParser(&self)
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	// get setting value
	settingRes, settingErr := settingGetAll(&Models.SettingSearch{})
	if settingErr != nil {
		return settingErr
	}
	byteArray, _ := json.Marshal(settingRes[0])
	var setting Models.Setting
	json.Unmarshal(byteArray, &setting)
	serialConvertValue := strconv.Itoa(setting.PurchaseSerial) // int to string to generate 9 digit code
	lenCheck := len(serialConvertValue)
	for i := 0; i < 9-lenCheck; i++ {
		serialConvertValue = "0" + serialConvertValue
	}
	self.Serial = serialConvertValue

	_, err = collection.InsertOne(context.Background(), self)

	if err == nil {
		// set setting value
		collectionSetting := DBManager.SystemCollections.Setting
		updateData := bson.M{
			"$set": bson.M{
				"purchaseserial": setting.PurchaseSerial + 1,
			},
		}
		_, updateErr := collectionSetting.UpdateOne(context.Background(), bson.M{"_id": setting.ID}, updateData)
		if updateErr != nil {
			c.Status(500)
			return errors.New("an error occurred when Incrementing Purchase Serial Number")
		} else {
			c.Status(200).Send([]byte(" Added New Purchase Successfully"))
			return nil
		}
	} else {
		c.Status(500)
		return err
	}
}

func PurchasingGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Purchasing

	// Fill the received search obj data
	var self Models.PurchasingSearch
	c.BodyParser(&self)

	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetPurchasingSearchBSONObj())
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

func PurchasingGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Purchasing
	var self Models.PurchasingSearch
	c.BodyParser(&self)
	b, results := Utils.FindByFilter(collection, self.GetPurchasingSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Purchasing
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]Models.PurchasingPopulated, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i], _ = PurchasingGetByIdPopulated(v.ID, &v)
	}
	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func PurchasingGetByIdPopulated(objID primitive.ObjectID, ptr *Models.Purchasing) (Models.PurchasingPopulated, error) {
	var PurchasingDoc Models.Purchasing
	if ptr == nil {
		PurchasingDoc, _ = PurchasingGetById(objID)
	} else {
		PurchasingDoc = *ptr
	}

	populatedResult := Models.PurchasingPopulated{}
	populatedResult.CloneFrom(PurchasingDoc)

	var err error

	// populate for SupplierRef
	if PurchasingDoc.SupplierRef != primitive.NilObjectID {
		populatedResult.SupplierRef, err = ContactGetById(PurchasingDoc.SupplierRef)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate for InventoryRef
	if PurchasingDoc.InventoryRef != primitive.NilObjectID {
		populatedResult.InventoryRef, err = InventoryGetById(PurchasingDoc.InventoryRef)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate for PurchaseOrderRef
	if PurchasingDoc.PurchaseOrderRef != primitive.NilObjectID {
		populatedResult.PurchaseOrderRef, err = PurchasingGetById(PurchasingDoc.PurchaseOrderRef)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate for Products of Invoice row array
	populatedResult.Products = make([]Models.InvoiceRowPurchasingPopulated, len(PurchasingDoc.Products))

	for i := range PurchasingDoc.Products {
		populatedResult.Products[i].CloneFrom(PurchasingDoc.Products[i])
		populatedResult.Products[i].ProductRef, err = ProductGetById(PurchasingDoc.Products[i].ProductRef)
		if err != nil {
			return populatedResult, err
		}
	}

	return populatedResult, nil
}

func PurchasingGetById(id primitive.ObjectID) (Models.Purchasing, error) {
	collection := DBManager.SystemCollections.Purchasing
	filter := bson.M{"_id": id}
	var self Models.Purchasing
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

// for modify Purchasing
func PurchasingModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Purchasing
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	_, err := PurchasingGetById(objID)
	if err != nil {
		return err
	}
	var self Models.Purchasing
	c.BodyParser(&self)
	err = self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	// need to add validate here for product using loop
	for _, product := range self.Products {
		err := product.Validate()
		if err != nil {
			return err
		}
	}

	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifying Purchasing")
	}

	c.Status(200).Send([]byte("Modified Successfully"))
	return nil
}

func PurchasingSetStatus(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Purchasing
	if c.Params("id") == "" || c.Params("new_status") == "" || c.Params("type") == "" {
		c.Status(404)
		return errors.New("all params not sent correctly")
	}

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	purchasingDB, _ := PurchasingGetById(objID) // need to add something
	newValue := purchasingDB.Status

	if c.Params("type") != purchasingDB.Type {
		c.Status(500)
		return errors.New("type doesn't match")
	}

	// for purchaseOrder
	if c.Params("type") == "purchaseorder" {
		if purchasingDB.Status == "created" && (c.Params("new_status") == "sent" || c.Params("new_status") == "decline") {
			newValue = c.Params("new_status")

		} else if purchasingDB.Status == "sent" && (c.Params("new_status") == "confirmed" || c.Params("new_status") == "decline") {
			newValue = c.Params("new_status")

		} else {
			c.Status(500)

			return errors.New("can't change status")
		}
	}

	// for purchaseDelivery
	if c.Params("type") == "purchasedelivery" {
		if purchasingDB.Status == "created" && (c.Params("new_status") == "chiped" || c.Params("new_status") == "decline") {
			newValue = c.Params("new_status")
		} else {
			c.Status(500)
			return errors.New("can't change status")
		}
	}

	updateData := bson.M{
		"$set": bson.M{
			"status": newValue,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifying purchasing status")
	} else {
		c.Status(200).Send([]byte("Modified Successfully"))
		return nil
	}
}

func SetConvertedToDelivery(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Purchasing
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	purchaseObj, err := PurchasingGetById(objID)
	if err != nil {
		return err
	}

	if purchaseObj.ConvertedToDelivery == false {
		purchaseObj.ConvertedToDelivery = true
	} else {
		return errors.New("Already converted")
	}

	updateData := bson.M{
		"$set": bson.M{
			"convertedtodelivery": purchaseObj.ConvertedToDelivery,
		},
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifying convertion")
	}

	c.Status(200).Send([]byte("Modified Successfully"))
	return nil
}
