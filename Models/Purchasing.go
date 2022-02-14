package Models

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchasing struct {
	ID                       primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	PurchaseOrderExternalRef string                 `json:"purchaseorderexternalref,omitempty"`
	Serial                   string                 `json:"serial,omitempty"`
	SupplierRef              primitive.ObjectID     `json:"supplierref,omitempty"` // populated done
	Date                     primitive.DateTime     `json:"date,omitempty"`
	Time                     primitive.DateTime     `json:"time,omitempty"`
	Status                   string                 `json:"status,omitempty"`
	ExpirationPeriod         int64                  `json:"expirationperiod,omitempty"`
	Products                 []InvoiceRowPurchasing `json:"products,omitempty"` // embedded populated
	TotalProduct             float64                `json:"totalproduct,omitempty"`
	TotalTax                 float64                `json:"totaltax,omitempty"`
	Total                    float64                `json:"total,omitempty"`
	InventoryRef             primitive.ObjectID     `json:"inventoryref,omitempty"` // for order  // populated done
	Type                     string                 `json:"type,omitempty"`
	PurchaseOrderRef         primitive.ObjectID     `json:"purchaseorderref,omitempty"` // populated done
	DeliveryCost             float64                `json:"deliverycost,omitempty"`
	DeliveryTicketId         string                 `json:"deliveryticketid,omitempty"`
	ConvertedToDelivery      bool                   `json:"convertedtodelivery,omitempty"`
	HasProducts              bool                   `json:"hasproducts,omitempty"`
}

type InvoiceRowPurchasing struct {
	ProductRef primitive.ObjectID `json:"productref,omitempty"`
	Amount     float64            `json:"amount,omitempty"` // same as demand
	Price      float64            `json:"price,omitempty"`
	Delivered  int                `json:"delivered,omitempty"`
}

func (obj Purchasing) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Type, validation.Required),
		validation.Field(&obj.Status, validation.Required),
	)
}

func (obj InvoiceRowPurchasing) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Amount, validation.Required),
		validation.Field(&obj.Price, validation.Required),
	)
}

type InvoiceRowPurchasingPopulated struct {
	ProductRef Product `json:"productref,omitempty"`
	Amount     float64 `json:"amount,omitempty"` // same as demand
	Price      float64 `json:"price,omitempty"`
	Delivered  int     `json:"delivered,omitempty"`
}

func (obj *InvoiceRowPurchasingPopulated) CloneFrom(other InvoiceRowPurchasing) {
	obj.ProductRef = Product{}
	obj.Amount = other.Amount
	obj.Price = other.Price
	obj.Delivered = other.Delivered
}

type PurchasingPopulated struct {
	ID                       primitive.ObjectID              `json:"_id,omitempty" bson:"_id,omitempty"`
	Serial                   string                          `json:"serial,omitempty"`
	PurchaseOrderExternalRef string                          `json:"purchaseorderexternalref,omitempty"`
	SupplierRef              Contact                         `json:"supplierref,omitempty"`
	Date                     primitive.DateTime              `json:"date,omitempty"`
	Time                     primitive.DateTime              `json:"time,omitempty"`
	Status                   string                          `json:"status,omitempty"`
	ExpirationPeriod         int64                           `json:"expirationperiod,omitempty"`
	Products                 []InvoiceRowPurchasingPopulated `json:"products,omitempty"`
	TotalProduct             float64                         `json:"totalproduct,omitempty"`
	TotalTax                 float64                         `json:"totaltax,omitempty"`
	Total                    float64                         `json:"total,omitempty"`
	InventoryRef             Inventory                       `json:"inventoryref,omitempty"` // for order
	Type                     string                          `json:"type,omitempty"`
	PurchaseOrderRef         Purchasing                      `json:"purchaseorderref,omitempty"`
	DeliveryCost             float64                         `json:"deliverycost,omitempty"`
	DeliveryTicketId         string                          `json:"deliveryticketid,omitempty"`
	ConvertedToDelivery      bool                            `json:"convertedtodelivery,omitempty"`
	HasProducts              bool                            `json:"hasproducts,omitempty"`
}

func (obj *PurchasingPopulated) CloneFrom(other Purchasing) {
	obj.ID = other.ID
	obj.Serial = other.Serial
	obj.SupplierRef = Contact{}
	obj.PurchaseOrderExternalRef = other.PurchaseOrderExternalRef
	obj.Date = other.Date
	obj.Time = other.Time
	obj.Status = other.Status
	obj.ExpirationPeriod = other.ExpirationPeriod
	obj.Products = []InvoiceRowPurchasingPopulated{}
	obj.TotalProduct = other.TotalProduct
	obj.TotalTax = other.TotalTax
	obj.Total = other.Total
	obj.InventoryRef = Inventory{}
	obj.Type = other.Type
	obj.PurchaseOrderRef = Purchasing{}
	obj.Serial = other.Serial
	obj.DeliveryCost = other.DeliveryCost
	obj.DeliveryTicketId = other.DeliveryTicketId
	obj.ConvertedToDelivery = other.ConvertedToDelivery
	obj.HasProducts = other.HasProducts
}

type PurchasingSearch struct {
	Type               string             `json:"type,omitempty"`
	TypeIsUsed         bool               `json:"typeisused,omitempty"`
	Status             string             `json:"status,omitempty"`
	StatusIsUsed       bool               `json:"statusisused,omitempty"`
	SupplierRef        primitive.ObjectID `json:"supplierref,omitempty"`
	SupplierRefIsUsed  bool               `json:"supplierrefisused,omitempty"`
	InventoryRef       primitive.ObjectID `json:"inventoryref,omitempty"`
	InventoryRefIsUsed bool               `json:"inventoryrefisused,omitempty"`
	DateFrom           primitive.DateTime `json:"datefrom,omitempty"`
	DateTo             primitive.DateTime `json:"dateto,omitempty"`
	DateRangeIsUsed    bool               `json:"daterangeisused,omitempty"`
	TimeFrom           primitive.DateTime `json:"timefrom,omitempty"`
	TimeTo             primitive.DateTime `json:"timeto,omitempty"`
	TimeRangeIsUsed    bool               `json:"timerangeisused,omitempty"`
}

func (obj PurchasingSearch) GetPurchasingSearchBSONObj() bson.M {
	self := bson.M{}

	if obj.TypeIsUsed {
		self["type"] = obj.Type
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.SupplierRefIsUsed {
		self["supplierref"] = obj.SupplierRef
	}

	if obj.InventoryRefIsUsed {
		self["inventoryref"] = obj.InventoryRef
	}

	if obj.DateRangeIsUsed {
		self["date"] = bson.M{
			"$gte": obj.DateFrom,
			"$lte": obj.DateTo,
		}
	}

	if obj.TimeRangeIsUsed {
		self["time"] = bson.M{
			"$gte": obj.TimeFrom,
			"$lte": obj.TimeTo,
		}
	}

	return self
}

func (obj Purchasing) GetModifcationBSONObj() bson.M {
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
