package Models

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"
	"errors"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DemandContent struct {
	ProductRef  primitive.ObjectID `json:"productref,omitempty" bson:"productref,omitempty"`
	MaterialRef primitive.ObjectID `json:"materialref,omitempty" bson:"materialref,omitempty"`
	Amount      float64            `json:"amount,omitempty" bson:"amount,omitempty"`
	Sent        float64            `json:"sent,omitempty" bson:"sent,omitempty"`
	Recieved    float64            `json:"recieved,omitempty" bson:"recieved,omitempty"`
}

type DemandContentPopulated struct {
	ProductRef  Product  `json:"productref,omitempty" bson:"productref,omitempty"`
	MaterialRef Material `json:"materialref,omitempty" bson:"materialref,omitempty"`
	Amount      float64  `json:"amount,omitempty" bson:"amount,omitempty"`
	Sent        float64  `json:"sent,omitempty" bson:"sent,omitempty"`
	Recieved    float64  `json:"recieved,omitempty" bson:"recieved,omitempty"`
}

type InventoryDemand struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FromInventory primitive.ObjectID `json:"frominventory,omitempty" bson:"frominventory,omitempty"`
	ToInventory   primitive.ObjectID `json:"toinventory,omitempty" bson:"toinventory,omitempty"`
	// UserRef       primitive.ObjectID `json:"userref"`
	DateTime      primitive.DateTime `json:"datetime,omitempty" bson:"datetime,omitempty"`
	Type          string             `json:"type,omitempty" bson:"type,omitempty"`
	RequestStatus bool               `json:"requeststatus,omitempty" bson:"requeststatus,omitempty"`
	EndDate       primitive.DateTime `json:"enddate,omitempty"`
	DemandList    []DemandContent    `json:"demandlist,omitempty" bson:"demandlist,omitempty"`
}

func (obj InventoryDemand) Validate() error {
	if obj.FromInventory == obj.ToInventory {
		return errors.New("to and from inventory cannot be same")
	}
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.ID, validation.Required),
		validation.Field(&obj.FromInventory, validation.Required),
		validation.Field(&obj.ToInventory, validation.Required),
		// validation.Field(&obj.UserRef, validation.Required),
		validation.Field(&obj.DemandList, validation.Required),
		validation.Field(&obj.Type, validation.Required, validation.In("Request", "Order", "Delivery")),
	)
}

type InventoryDemandPopulated struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FromInventory Inventory          `json:"frominventory,omitempty" bson:"frominventory,omitempty"`
	ToInventory   Inventory          `json:"toinventory,omitempty" bson:"toinventory,omitempty"`
	// UserRef       User `json:"userref"`
	DateTime      primitive.DateTime       `json:"datetime,omitempty" bson:"datetime,omitempty"`
	Type          string                   `json:"type,omitempty" bson:"type,omitempty"`
	RequestStatus bool                     `json:"requeststatus,omitempty" bson:"requeststatus,omitempty"`
	DemandList    []DemandContentPopulated `json:"demandlist,omitempty" bson:"demandlist,omitempty"`
	EndDate       primitive.DateTime       `json:"enddate,omitempty"`
}

func (obj *InventoryDemandPopulated) CloneFrom(other InventoryDemand) {
	obj.ID = other.ID
	obj.ToInventory = Inventory{}
	obj.FromInventory = Inventory{}
	// obj.UserRef = User{}
	obj.DateTime = other.DateTime
	obj.Type = other.Type
	obj.RequestStatus = other.RequestStatus
	obj.DemandList = []DemandContentPopulated{}
	obj.EndDate = other.EndDate
}

type InventoryDemandSearch struct {
	IDIsUsed            bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FromInventoryIsUsed bool               `json:"frominventoryisused,omitempty" bson:"frominventoryisused,omitempty"`
	FromInventory       primitive.ObjectID `json:"frominventory,omitempty" bson:"frominventory,omitempty"`
	ToInventoryIsUsed   bool               `json:"toinventoryisused,omitempty" bson:"toinventoryisused,omitempty"`
	ToInventory         primitive.ObjectID `json:"toinventory,omitempty" bson:"toinventory,omitempty"`
	// UserRefIsUsed bool `json:"userrefisused:omitempty"`
	// UserRef       primitive.ObjectID `json:"userref"`
	DateTimeIsUsed      bool               `json:"datetimeisused,omitempty" bson:"datetimeisused,omitempty"`
	DateTime            primitive.DateTime `json:"datetime,omitempty" bson:"datetime,omitempty"`
	TypeIsUsed          bool               `json:"typeisused,omitempty" bson:"typeisused,omitmepty"`
	Type                string             `json:"type,omitempty" bson:"type,omitempty"`
	RequestStatusIsUsed bool               `json:"requeststatusisused,omitempty" bson:"requeststatusisused,omitmepty"`
	RequestStatus       bool               `json:"requeststatus,omitempty" bson:"requeststatus,omitempty"`
}

func (obj InventoryDemand) GetModifcationBSONObj() bson.M {
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

func (obj InventoryDemandSearch) GetInventoryDemandSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	// if obj.UserRefIsUsed{
	// 	self['userref']=obj.UserRef
	// }

	if obj.FromInventoryIsUsed {
		self["frominventory"] = obj.FromInventory
	}

	if obj.ToInventoryIsUsed {
		self["toinventory"] = obj.ToInventory
	}

	if obj.DateTimeIsUsed {
		self["datetime"] = obj.DateTime
	}

	if obj.TypeIsUsed {
		self["type"] = obj.Type
	}

	if obj.RequestStatusIsUsed {
		self["requeststatus"] = obj.RequestStatus
	}

	return self
}
