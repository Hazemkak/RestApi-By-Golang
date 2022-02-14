package Models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UOMConvert struct {
	FromUnitRef primitive.ObjectID `json:"fromunitref,omitempty" bson:"fromunitref,omitempty"`
	ToUnitRef   primitive.ObjectID `json:"tounitref,omitempty" bson:"tounitref,omitempty"`
	Amount      float64            `json:"amount,omitempty"`
}

type UOMRelation struct {
	UnitRef primitive.ObjectID `json:"unitref,omitempty" bson:"unitref,omitempty"`
	Ratio   float64            `json:"ratio,omitempty"`
	Status  bool               `json:"status,omitempty"`
}

type UOMRelationPopulated struct {
	UnitRef UnitOfMeasurement `json:"unitref,omitempty" bson:"unitref,omitempty"`
	Ratio   float64           `json:"ratio,omitempty"`
	Status  bool              `json:"status"`
}

type UnitOfMeasurement struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Status    bool               `json:"status,omitempty"`
	Category  string             `json:"category,omitempty"`
	Relations []UOMRelation      `json:"relations,omitempty"`
}

type UnitOfMeasurementPopulated struct {
	ID        primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Status    bool                   `json:"status"`
	Category  string                 `json:"category,omitempty"`
	Relations []UOMRelationPopulated `json:"relations,omitempty"`
}

type UnitOfMeasurementSearch struct {
	ID             string `json:"_id" bson:"_id"`
	IDIsUsed       bool   `json:"idisused"`
	Name           string `json:"name,omitempty"`
	NameIsUsed     bool   `json:"nameisused,omitempty"`
	Status         bool   `json:"status,omitempty"`
	StatusIsUsed   bool   `json:"statusisused,omitempty"`
	Category       string `json:"category,omitempty"`
	CategoryIsUsed bool   `json:"categoryisused,omitempty"`
}

func (obj UnitOfMeasurement) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
	)
}

// https://docs.mongodb.com/manual/reference/operator/update/push/#mongodb-update-up.-push
func UOMRelationGetAppendBSONObj(arr []UOMRelation) bson.M {
	each := bson.M{
		"$each": []bson.M{},
	}

	// https://stackoverflow.com/questions/50939497/golang-cast-interface-to-struct
	for _, v := range arr {
		each["$each"] = append(each["$each"].([]bson.M), bson.M{
			"unitref": v.UnitRef,
			"ratio":   v.Ratio,
			"status":  v.Status,
		})
	}

	self := bson.M{
		"relations": each,
	}

	return self
}

func (obj UnitOfMeasurementSearch) GetUnitOfMeasurementSearchBSONObj() bson.M {
	self := bson.M{}

	if obj.IDIsUsed {
		self["_id"], _ = primitive.ObjectIDFromHex(obj.ID)
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.CategoryIsUsed {
		self["category"] = obj.Category
	}

	return self
}

func (obj *UnitOfMeasurementPopulated) CloneFrom(other UnitOfMeasurement) {
	obj.ID = other.ID
	obj.Relations = []UOMRelationPopulated{}
	obj.Category = other.Category
	obj.Name = other.Name
	obj.Status = other.Status
}
