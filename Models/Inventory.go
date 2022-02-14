package Models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inventory struct {
	ID            primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	InventoryName string                `json:"inventoryname,omitempty"`
	Type          string                `json:"type,omitempty"`
	Address       string                `json:"address,omitempty"`
	Longitude     float64               `json:"longitude,omitempty"`
	Latitude      float64               `json:"latitude,omitempty"`
	Status        bool                  `json:"status,omitempty"`
	Contents      []InventoryContentRow `json:"contents,omitempty"`
}

type InventoryContentRow struct {
	ProductRef     primitive.ObjectID `json:"productref,omitempty"`
	RawMaterialRef primitive.ObjectID `json:"rawmaterialref,omitempty"`
	Amount         float64            `json:"amount,omitempty"`
}

type InventoryPopulated struct {
	ID            primitive.ObjectID             `json:"_id,omitempty" bson:"_id,omitempty"`
	InventoryName string                         `json:"inventoryname,omitempty"`
	Type          string                         `json:"type,omitempty"`
	Address       string                         `json:"address,omitempty"`
	Longitude     float64                        `json:"longitude,omitempty"`
	Latitude      float64                        `json:"latitude,omitempty"`
	Status        bool                           `json:"status,omitempty"`
	Contents      []InventoryContentRowPopulated `json:"contents,omitempty"`
}

type InventoryContentRowPopulated struct {
	ProductRef     Product  `json:"productref,omitempty"`
	RawMaterialRef Material `json:"rawmaterialref,omitempty"`
	Amount         float64  `json:"amount,omitempty"`
}

func (obj Inventory) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.InventoryName, validation.Required),
		validation.Field(&obj.Type, validation.In("Product Inventory", "Raw Material Inventory")),
	)
}

func (obj InventoryContentRow) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Amount, validation.Required),
	)
}

func (obj *InventoryContentRowPopulated) CloneFrom(other InventoryContentRow) {
	obj.ProductRef = Product{}
	obj.RawMaterialRef = Material{}
	obj.Amount = other.Amount
}

func (obj *InventoryPopulated) CloneFrom(other Inventory) {
	obj.ID = other.ID
	obj.InventoryName = other.InventoryName
	obj.Type = other.Type
	obj.Address = other.Address
	obj.Longitude = other.Longitude
	obj.Latitude = other.Latitude
	obj.Status = other.Status
	obj.Contents = []InventoryContentRowPopulated{}
}

type InventorySearch struct {
	InventoryName       string `json:"inventoryname,omitempty"`
	InventoryNameIsUsed bool   `json:"inventorynameisused,omitempty"`
	Type                string `json:"type,omitempty"`
	TypeIsUsed          bool   `json:"typeisused,omitempty"`
	Status              bool   `json:"status,omitempty"`
	StatusIsUsed        bool   `json:"statusisused,omitempty"`
}

func (obj InventorySearch) GetBSONSearchObj() bson.M {
	self := bson.M{}
	if obj.InventoryNameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.InventoryName)
		self["inventoryname"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}
	if obj.TypeIsUsed {
		self["type"] = obj.Type
	}
	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}
	return self
}

func (obj Inventory) GetBSONModificationObj() bson.M {
	self := bson.M{
		"inventoryname": obj.InventoryName,
		"type":          obj.Type,
		"address":       obj.Address,
		"longitude":     obj.Longitude,
		"latitude":      obj.Latitude,
		"status":        obj.Status,
		"contents":      obj.Contents,
	}
	return self
}
