package Models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Material struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name,omitempty"`
	Price          float64            `json:"price,omitempty"`
	Status         bool               `json:"status,omitempty"`
	Diameter       float64            `json:"diameter,omitempty"`
	DiameterUomId  primitive.ObjectID `json:"diameteruomid,omitempty" bson:"diameteruomid,omitempty"`
	Weight         float64            `json:"weight,omitempty"`
	WeightUomId    primitive.ObjectID `json:"weightuomid,omitempty" bson:"weightuomid,omitempty"`
	Length         float64            `json:"length,omitempty"`
	LengthUomId    primitive.ObjectID `json:"lengthuomid,omitempty" bson:"lengthuomid,omitempty"`
	Thickness      float64            `json:"thickness,omitempty"`
	ThicknessUomId primitive.ObjectID `json:"thicknessuomid,omitempty" bson:"thicknessuomid,omitempty"`
}

func (obj Material) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
	)
}

type MaterialPopulated struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name,omitempty"`
	Price          float64            `json:"price,omitempty"`
	Status         bool               `json:"status,omitempty"`
	Diameter       float64            `json:"diameter,omitempty"`
	DiameterUomId  UnitOfMeasurement  `json:"diameteruomid,omitempty" bson:"diagramuomid,omitempty"`
	Weight         float64            `json:"weight,omitempty"`
	WeightUomId    UnitOfMeasurement  `json:"weightuomid,omitempty" bson:"weightuomid,omitempty"`
	Length         float64            `json:"length,omitempty"`
	LengthUomId    UnitOfMeasurement  `json:"lengthuomid,omitempty" bson:"lengthuomid,omitempty"`
	Thickness      float64            `json:"thickness,omitempty"`
	ThicknessUomId UnitOfMeasurement  `json:"thicknessuomid,omitempty" bson:"thicknessuomid,omitempty"`
}

func (obj *MaterialPopulated) CloneFrom(other Material) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.Price = other.Price
	obj.Status = other.Status
	obj.Diameter = other.Diameter
	obj.DiameterUomId = UnitOfMeasurement{}
	obj.Weight = other.Weight
	obj.WeightUomId = UnitOfMeasurement{}
	obj.Length = other.Length
	obj.LengthUomId = UnitOfMeasurement{}
	obj.Thickness = other.Thickness
	obj.ThicknessUomId = UnitOfMeasurement{}
}

func (obj Material) GetBSONModificationObj() bson.M {
	self := bson.M{
		"name":           obj.Name,
		"price":          obj.Price,
		"status":         obj.Status,
		"diameter":       obj.Diameter,
		"diameteruomid":  obj.DiameterUomId,
		"weight":         obj.Weight,
		"weightuomid":    obj.WeightUomId,
		"length":         obj.Length,
		"lengthuomid":    obj.LengthUomId,
		"thickness":      obj.Thickness,
		"thicknessuomid": obj.ThicknessUomId,
	}
	return self
}

type MaterialSearch struct {
	Name            string `json:"name,omitempty"`
	NameIsUsed      bool   `json:"nameisused,omitempty"`
	Price           string `json:"price,omitempty"`
	PriceIsUsed     bool   `json:"priceisused,omitempty"`
	Status          bool   `json:"status,omitempty"`
	StatusIsUsed    bool   `json:"statusisused,omitempty"`
	Diameter        string `json:"diameter,omitempty"`
	DiameterIsUsed  bool   `json:"diameterisused,omitempty"`
	Weight          string `json:"weight,omitempty"`
	WeightIsUsed    bool   `json:"weightisused,omitempty"`
	Length          string `json:"length,omitempty"`
	LengthIsUsed    bool   `json:"lengthisused,omitempty"`
	Thickness       string `json:"thickness,omitempty"`
	ThicknessIsUsed bool   `json:"thicknessisused,omitempty"`
}

func (obj MaterialSearch) GetMaterialSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.PriceIsUsed {
		self["price"] = obj.Price
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.DiameterIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Diameter)
		self["diameter"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.WeightIsUsed {
		self["weigth"] = obj.Weight
	}

	if obj.LengthIsUsed {
		self["length"] = obj.Length
	}

	if obj.ThicknessIsUsed {
		self["thickness"] = obj.Thickness
	}

	return self
}
