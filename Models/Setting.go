package Models

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Setting struct {
	ID                  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductSerial       int                `json:"productserial,omitempty"`
	SalesSerial         int                `json:"salesserial,omitempty"`
	PurchaseSerial      int                `json:"purchaseserial,omitempty"`
	SalesDeliverySerial int                `json:"salesdeliveryserial,omitempty"`
}

type SettingSearch struct {
	IDIsUsed                  bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID                        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductSerial             int                `json:"productserial,omitempty"`
	ProductSerialIsUsed       bool               `json:"productserialisused,omitempty"`
	SalesSerial               int                `json:"salesserial,omitempty"`
	SalesSerialIsUsed         bool               `json:"salesserialisused,omitempty"`
	PurchaseSerial            int                `json:"purchaseserial,omitempty"`
	PurchaseSerialIsUsed      bool               `json:"purchaseserialisused,omitempty"`
	SalesDeliverySerial       int                `json:"salesdeliveryserial,omitempty"`
	SalesDeliverySerialIsUsed bool               `json:"salesdeliveryserialisused,omitempty`
}

func (obj Setting) GetIdString() string {
	return obj.ID.String()
}

func (obj Setting) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj SettingSearch) GetSettingSearchBSONObj() bson.M {
	self := bson.M{}
	if obj.IDIsUsed {
		self["_id"] = obj.ID
	}

	if obj.ProductSerialIsUsed {
		self["productserial"] = obj.ProductSerial
	}
	if obj.SalesSerialIsUsed {
		self["salesserial"] = obj.SalesSerial
	}
	if obj.PurchaseSerialIsUsed {
		self["purchaseserial"] = obj.PurchaseSerial
	}

	if obj.SalesSerialIsUsed {
		self["salesserial"] = obj.SalesSerial
	}
	if obj.PurchaseSerialIsUsed {
		self["purchaseserial"] = obj.PurchaseSerial
	}

	if obj.SalesDeliverySerialIsUsed {
		self["salesdeliveryserial"] = obj.SalesDeliverySerial
	}

	return self
}
