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

type ProductOffer struct {
	ProductRef primitive.ObjectID `json:"productref,omitempty" bson:"productref,omitempty"`
	OfferPrice float64            `json:"offerprice,omitempty"`
}

type ProductOfferPopulated struct {
	ProductRef Product `json:"productref,omitempty" bson:"productref,omitempty"`
	OfferPrice float64 `json:"offerprice,omitempty"`
}

type PriceList struct {
	ID                 primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name               string               `json:"name"`
	Status             bool                 `json:"status"`
	AppliedForContacts []primitive.ObjectID `json:"appliedforcontacts,omitempty" bson:"appliedforcontacts,omitempty"`
	ProductsList       []ProductOffer       `json:"productslist"`
}

func (obj PriceList) GetIdString() string {
	return obj.ID.String()
}

func (obj PriceList) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj PriceList) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.ProductsList, validation.Required),
	)
}

func (obj PriceList) GetModifcationBSONObj() bson.M {
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

type PriceListSearch struct {
	IDIsUsed     bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NameIsUsed   bool               `json:"nameisused,omitempty"`
	Name         string             `json:"name,omitempty"`
	StatusIsUsed bool               `json:"statusisused,omitempty"`
	Status       bool               `json:"status,omitempty"`
}

func (obj PriceListSearch) GetPriceListSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.NameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Name)
		self["name"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	return self
}

type PriceListPopulated struct {
	ID                 primitive.ObjectID      `json:"_id,omitempty" bson:"_id,omitempty"`
	Name               string                  `json:"name,omitempty"`
	Status             bool                    `json:"status,omitempty"`
	AppliedForContacts []Contact               `json:"appliedforcontacts,omitempty" bson:"appliedforcontacts,omitempty"`
	ProductsList       []ProductOfferPopulated `json:"productslist,omitempty"`
}

func (obj *PriceListPopulated) CloneFrom(other PriceList) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.Status = other.Status
	obj.AppliedForContacts = []Contact{}
	obj.ProductsList = []ProductOfferPopulated{}
}

func (obj *ProductOfferPopulated) CloneFrom(other ProductOffer) {
	obj.ProductRef = Product{}
	obj.OfferPrice = other.OfferPrice
}
