package Models

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"
	"fmt"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name              string             `json:"name,omitempty"`
	Type              string             `json:"type,omitempty"`
	Serial            string             `json:"serial"`
	Status            bool               `json:"status,omitempty"`
	CanBeSold         bool               `json:"canbesold,omitempty"`
	CanBePurchased    bool               `json:"canbepurchased,omitempty"`
	CanBeManufactured bool               `json:"canbemanufactured,omitempty"`
	Category          string             `json:"category,omitempty"`
	Barcode           string             `json:"barcode,omitempty"`
	ExternalReference string             `json:"externalreference,omitempty"`
	SalesPrice        float64            `json:"salesprice,omitempty"`
	CostPrice         float64            `json:"costprice,omitempty"`
	PurchasedUomId    primitive.ObjectID `json:"purchaseduomid,omitempty" bson:"purchaseduomid,omitempty"`
	SalesUomId        primitive.ObjectID `json:"salesuomid,omitempty" bson:"salesuomid,omitempty"`
	ManufacturedUomId primitive.ObjectID `json:"manufactureduomid,omitempty" bson:"manufactureduomid,omitempty"`
	Weight            float64            `json:"weight,omitempty"`
	WeightUomId       primitive.ObjectID `json:"weightuomid,omitempty" bson:"weightuomid,omitempty"`
	Volume            float64            `json:"volume,omitempty"`
	VolumeUomId       primitive.ObjectID `json:"volumeuomid,omitempty" bson:"volumeuomid,omitempty"`
	Length            float64            `json:"length,omitempty"`
	LengthUomId       primitive.ObjectID `json:"lengthuomid,omitempty" bson:"lengthuomid,omitempty"`
	DeliveryNote      string             `json:"deliverynote,omitempty"`
	ReceiptNote       string             `json:"receiptnote,omitempty"`
	BillOfMaterial    []BOM              `json:"billofmaterial"`
}

func (obj Product) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Type, validation.Required),
		validation.Field(&obj.Category, validation.Required),
		validation.Field(&obj.Barcode, validation.Required),
		validation.Field(&obj.ExternalReference, validation.Required),
		validation.Field(&obj.SalesPrice, validation.Required),
		validation.Field(&obj.CostPrice, validation.Required),
		validation.Field(&obj.Weight, validation.Required),
		validation.Field(&obj.Volume, validation.Required),
		validation.Field(&obj.Length, validation.Required),
	)
}

type BOM struct {
	RawMaterialRef primitive.ObjectID `json:"rawmaterialref,omitempty" bson:"rawmaterialref,omitempty"`
	Quantity       float64            `json:"quantity,omitempty"`
	Status         bool               `json:"status,omitempty"`
}

func (obj BOM) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Quantity, validation.Required),
		validation.Field(&obj.RawMaterialRef, validation.Required),
	)
}

type BOMPopulated struct {
	RawMaterialRef Material `json:"rawMaterialRef,omitempty"`
	Quantity       float64  `json:"quantity,omitempty"`
	Status         bool     `json:"status,omitempty"`
}

func (obj *BOMPopulated) CloneFrom(other BOM) {
	obj.RawMaterialRef = Material{}
	obj.Quantity = other.Quantity
	obj.Status = other.Status
}

type ProductPopulated struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name              string             `json:"name,omitempty"`
	Type              string             `json:"type,omitempty"`
	Serial            string             `json:"serial"`
	Status            bool               `json:"status,omitempty"`
	CanBeSold         bool               `json:"canbesold,omitempty"`
	CanBePurchased    bool               `json:"canbepurchased,omitempty"`
	CanBeManufactured bool               `json:"canbemanufactured,omitempty"`
	Category          string             `json:"category,omitempty"`
	Barcode           string             `json:"barcode,omitempty"`
	ExternalReference string             `json:"externalreference,omitempty"`
	SalesPrice        float64            `json:"salesprice,omitempty"`
	CostPrice         float64            `json:"costprice,omitempty"`
	PurchasedUomId    UnitOfMeasurement  `json:"purchaseduomid,omitempty" bson:"purchaseduomid,omitempty"`
	SalesUomId        UnitOfMeasurement  `json:"salesuomid,omitempty" bson:"salesuomid,omitempty"`
	ManufacturedUomId UnitOfMeasurement  `json:"manufactureduomid,omitempty" bson:"manufactureduomid,omitempty"`
	Weight            float64            `json:"weight,omitempty"`
	WeightUomId       UnitOfMeasurement  `json:"weightuomid,omitempty" bson:"weightuomid,omitempty"`
	Volume            float64            `json:"volume,omitempty"`
	VolumeUomId       UnitOfMeasurement  `json:"volumeuomid,omitempty" bson:"volumeuomid,omitempty"`
	Length            float64            `json:"length,omitempty"`
	LengthUomId       UnitOfMeasurement  `json:"lengthuomid,omitempty" bson:"lengthuomid,omitempty"`
	DeliveryNote      string             `json:"deliverynote,omitempty"`
	ReceiptNote       string             `json:"receiptnote,omitempty"`
	BillOfMaterial    []BOMPopulated     `json:"billofmaterial"`
}

func (obj *ProductPopulated) CloneFrom(other Product) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.Type = other.Type
	obj.Status = other.Status
	obj.Serial = other.Serial
	obj.CanBeSold = other.CanBeSold
	obj.CanBePurchased = other.CanBePurchased
	obj.CanBeManufactured = other.CanBeManufactured
	obj.Category = other.Category
	obj.Barcode = other.Barcode
	obj.ExternalReference = other.ExternalReference
	obj.SalesPrice = other.SalesPrice
	obj.CostPrice = other.CostPrice
	obj.PurchasedUomId = UnitOfMeasurement{}
	obj.SalesUomId = UnitOfMeasurement{}
	obj.ManufacturedUomId = UnitOfMeasurement{}
	obj.WeightUomId = UnitOfMeasurement{}
	obj.VolumeUomId = UnitOfMeasurement{}
	obj.LengthUomId = UnitOfMeasurement{}
	obj.Weight = other.Weight
	obj.Length = other.Length
	obj.Volume = other.Volume
	obj.DeliveryNote = other.DeliveryNote
	obj.ReceiptNote = other.ReceiptNote
	obj.BillOfMaterial = []BOMPopulated{}
}

type ProductSearch struct {
	Name                    string `json:"name,omitempty"`
	NameIsUsed              bool   `json:"nameisused,omitempty"`
	Category                string `json:"category,omitempty"`
	CategoryIsUsed          bool   `json:"categoryisused,omitempty"`
	Status                  bool   `json:"status,omitempty"`
	StatusIsUsed            bool   `json:"statusisused,omitempty"`
	Serial                  string `json:"serial,omitempty"`
	SerialIsUsed            bool   `json:"serialisused,omitempty"`
	CanBeSold               bool   `json:"canbesold,omitempty"`
	CanBeSoldIsUsed         bool   `json:"canbesoldisused,omitempty"`
	CanBePurchased          bool   `json:"purchased,omitempty"`
	CanBePurchasedIsUsed    bool   `json:"canbepurchasedisused,omitempty"`
	CanBeManufactured       bool   `json:"canbemanufactured,omitempty"`
	CanBeManufacturedIsUsed bool   `json:"canbemanufacturedisused,omitempty"`
}

func (obj Product) GetModifcationBSONObj() bson.M {
	self := bson.M{}
	valueOfObj := reflect.ValueOf(obj)
	typeOfObj := valueOfObj.Type()
	invalidFieldNames := []string{"ID"}

	for i := 0; i < valueOfObj.NumField(); i++ {
		if Utils.ArrayStringContains(invalidFieldNames, typeOfObj.Field(i).Name) {
			continue
		}
		self[strings.ToLower(typeOfObj.Field(i).Name)] = valueOfObj.Field(i).Interface()
	}
	return self
}

func (obj ProductSearch) GetProductSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.CategoryIsUsed {
		self["category"] = obj.Category
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.SerialIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Serial)
		self["serial"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}

	}

	if obj.CanBeSoldIsUsed {
		self["canbesold"] = obj.CanBeSold
	}

	if obj.CanBePurchasedIsUsed {
		self["canbepurchased"] = obj.CanBePurchased
	}

	if obj.CanBeManufacturedIsUsed {
		self["canbemanufactured"] = obj.CanBeManufactured
	}

	return self
}
