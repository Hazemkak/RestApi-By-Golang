package Controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"SEEN-TECH-VAI21-BACKEND-GO/DBManager"
	"SEEN-TECH-VAI21-BACKEND-GO/Models"
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func isProductExisting(name string) (bool, interface{}) {
	collection := DBManager.SystemCollections.Product
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

func ProductNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Product

	var self Models.Product
	c.BodyParser(&self)

	_, existing := isProductExisting(self.Name)
	if existing != "" {
		return errors.New("this Name already exists with same Product Name")
	}
	// Validate the obj
	err := self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	// get setting value
	res, err := settingGetAll(&Models.SettingSearch{})

	if err != nil {
		return err
	}
	byteArray, _ := json.Marshal(res[0])
	var setting Models.Setting
	json.Unmarshal(byteArray, &setting)
	serialConvertValue := strconv.Itoa(setting.ProductSerial) // int to string to generate hex code
	self.Serial = fmt.Sprintf("%09x", serialConvertValue)     // convert 8 digit number

	//added for make an array of BOM and check here for fronted value
	if len(self.BillOfMaterial) == 0 {
		self.BillOfMaterial = []Models.BOM{}
	}

	_, err = collection.InsertOne(context.Background(), self)
	if err == nil {

		// set setting value
		collectionSetting := DBManager.SystemCollections.Setting
		updateData := bson.M{
			"$set": bson.M{
				"productserial": setting.ProductSerial + 1,
			},
		}
		_, updateErr := collectionSetting.UpdateOne(context.Background(), bson.M{"_id": setting.ID}, updateData)
		if updateErr != nil {
			c.Status(500)
			return errors.New("an error occurred when Incrementing Serial Number")
		} else {
			c.Status(200).Send([]byte(" Added New Product Successfully"))
			return nil
		}

	} else {
		c.Status(500)
		return err
	}

}

func ProductGetAll(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Product

	// Fill the received search obj data
	var self Models.ProductSearch
	c.BodyParser(&self)

	var results []bson.M
	b, results := Utils.FindByFilter(collection, self.GetProductSearchBSONObj())
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

func ProductGetAllPopulated(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Product
	var self Models.ProductSearch
	c.BodyParser(&self)
	b, results := Utils.FindByFilter(collection, self.GetProductSearchBSONObj())
	if !b {
		c.Status(500)
		return errors.New("object is not found")
	}
	byteArr, _ := json.Marshal(results)
	var ResultDocs []Models.Product
	json.Unmarshal(byteArr, &ResultDocs)
	populatedResult := make([]Models.ProductPopulated, len(ResultDocs))

	for i, v := range ResultDocs {
		populatedResult[i], _ = ProductGetByIdPopulated(v.ID, &v)
	}
	allpopulated, _ := json.Marshal(bson.M{"result": populatedResult})
	c.Set("Content-Type", "application/json")
	c.Send(allpopulated)
	return nil
}

func ProductGetByIdPopulated(objID primitive.ObjectID, ptr *Models.Product) (Models.ProductPopulated, error) {
	var ProductDoc Models.Product
	if ptr == nil {
		ProductDoc, _ = ProductGetById(objID)
	} else {
		ProductDoc = *ptr
	}

	populatedResult := Models.ProductPopulated{}
	populatedResult.CloneFrom(ProductDoc)

	var err error

	// populate for purchaseduomid
	if ProductDoc.PurchasedUomId != primitive.NilObjectID {
		populatedResult.PurchasedUomId, err = UnitsOfMeasurementGetById(ProductDoc.PurchasedUomId)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate for salesuomid
	if ProductDoc.SalesUomId != primitive.NilObjectID {
		populatedResult.SalesUomId, err = UnitsOfMeasurementGetById(ProductDoc.SalesUomId)
		if err != nil {
			return populatedResult, err
		}
	}
	// populate for ManufacturedUomId
	if ProductDoc.ManufacturedUomId != primitive.NilObjectID {
		populatedResult.ManufacturedUomId, err = UnitsOfMeasurementGetById(ProductDoc.ManufacturedUomId)
		if err != nil {
			return populatedResult, err
		}
	}
	// populate for WeightUomId
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
	// populate for VolumeUomId
	if ProductDoc.VolumeUomId != primitive.NilObjectID {
		populatedResult.VolumeUomId, err = UnitsOfMeasurementGetById(ProductDoc.VolumeUomId)
		if err != nil {
			return populatedResult, err
		}
	}

	// populate for Raw Material of BOM array
	populatedResult.BillOfMaterial = make([]Models.BOMPopulated, len(ProductDoc.BillOfMaterial))
	for i := range ProductDoc.BillOfMaterial {
		populatedResult.BillOfMaterial[i].CloneFrom(ProductDoc.BillOfMaterial[i])
		populatedResult.BillOfMaterial[i].RawMaterialRef, err = MaterialGetById(ProductDoc.BillOfMaterial[i].RawMaterialRef)
		if err != nil {
			return populatedResult, err
		}
	}

	return populatedResult, nil
}

func ProductGetById(id primitive.ObjectID) (Models.Product, error) {
	collection := DBManager.SystemCollections.Product
	filter := bson.M{"_id": id}
	var self Models.Product
	_, results := Utils.FindByFilter(collection, filter)
	if len(results) == 0 {
		return self, errors.New("obj not found")
	}
	bsonBytes, _ := bson.Marshal(results[0]) // Decode
	bson.Unmarshal(bsonBytes, &self)         // Encode
	return self, nil
}

// for modify product
func ProductModify(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Product
	if c.Params("id") == "" {
		c.Status(404)
		return errors.New("id param needed")
	}
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))

	_, err := ProductGetById(objID)
	if err != nil {
		return err
	}
	var self Models.Product
	c.BodyParser(&self)
	err = self.Validate()
	if err != nil {
		c.Status(500)
		return err
	}
	// need to add validate here for BOM using loop
	for _, bom := range self.BillOfMaterial {
		err := bom.Validate()
		if err != nil {
			return err
		}
	}
	// Check if name already exists
	_, id := isProductExisting(self.Name)
	if id != "" && id != objID.Hex() {
		c.Status(500)
		return errors.New("name already exists")
	}

	updateData := bson.M{
		"$set": self.GetModifcationBSONObj(),
	}
	_, updateErr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, updateData)
	if updateErr != nil {
		c.Status(500)
		return errors.New("an error occurred when modifying product")
	}

	c.Status(200).Send([]byte("Modified Successfully"))
	return nil
}

func ProductGetDistinctCategories(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Product

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

func ProductSetStatus(c *fiber.Ctx) error {
	// Initiate the connection
	collection := DBManager.SystemCollections.Product

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
		return errors.New("an error occurred when modifying product status")
	}

	c.Status(200).Send([]byte("status modified successfully"))
	return nil
}

// for add new BOM in Product
func ProductAddBOMNew(c *fiber.Ctx) error {
	collection := DBManager.SystemCollections.Product
	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))

	var newBOM Models.BOM
	c.BodyParser(&newBOM)

	err := newBOM.Validate()
	if err != nil {
		c.Status(500)
		return err
	}

	filter := bson.M{"_id": objId}
	updateData := bson.M{
		"$push": bson.M{
			"billofmaterial": newBOM,
		},
	}

	_, updateErr := collection.UpdateOne(context.Background(), filter, updateData)
	if updateErr != nil {
		return errors.New("an error occurred when adding new BOM")
	}
	c.Status(200).Send([]byte("Added New BOM In the Product Successfully"))
	return nil
}
