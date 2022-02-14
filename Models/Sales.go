package Models

import (
	"SEEN-TECH-VAI21-BACKEND-GO/Utils"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sales struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Serial           string             `json:"serial,omitempty"`
	CustomerRef      primitive.ObjectID `json:"customerref,omitempty"` // populated done
	Date             primitive.DateTime `json:"date,omitempty"`
	Time             primitive.DateTime `json:"time,omitempty"`
	Status           string             `json:"status,omitempty"`
	PriceListRef     primitive.ObjectID `json:"pricelistref,omitempty"` // populated done
	ExpirationPeriod int64              `json:"expirationperiod,omitempty"`
	Products         []InvoiceRow       `json:"products,omitempty"` // embedded populated
	TotalProduct     float64            `json:"totalproduct,omitempty"`
	TotalTax         float64            `json:"totaltax,omitempty"`
	Total            float64            `json:"total,omitempty"`
	InventoryRef     primitive.ObjectID `json:"inventoryref,omitempty"` // for order  // populated done
	Type             string             `json:"type,omitempty"`
	QuotationRef     primitive.ObjectID `json:"quotationref,omitempty"`  // populated  done
	SalesOrderRef    primitive.ObjectID `json:"salesorderref,omitempty"` // populated done
	ConvertedToSales bool               `json:"convertedtosales,omitempty"`
	IsDelivered      bool               `json:"isdelivered,omitempty"`
}

type InvoiceRow struct {
	ProductRef primitive.ObjectID `json:"productref,omitempty"`
	Amount     float64            `json:"amount,omitempty"` // same as demand
	Price      float64            `json:"price,omitempty"`
	Delivered  int                `json:"delivered,omitempty"`
}

// need to change when fronted is ready
func (obj Sales) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Type, validation.Required),
		validation.Field(&obj.Status, validation.Required),
		validation.Field(&obj.CustomerRef, validation.Required),
	)
}

func (obj InvoiceRow) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Amount, validation.Required),
		validation.Field(&obj.Price, validation.Required),
	)
}

type InvoiceRowPopulated struct {
	ProductRef Product `json:"productref,omitempty"`
	Amount     float64 `json:"amount,omitempty"` // same as demand
	Price      float64 `json:"price,omitempty"`
	Delivered  int     `json:"delivered,omitempty"`
}

func (obj *InvoiceRowPopulated) CloneFrom(other InvoiceRow) {
	obj.ProductRef = Product{}
	obj.Amount = other.Amount
	obj.Price = other.Price
	obj.Delivered = other.Delivered
}

type SalesPopulated struct {
	ID               primitive.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Serial           string                `json:"serial,omitempty"`
	CustomerRef      Contact               `json:"customerref,omitempty"`
	Date             primitive.DateTime    `json:"date,omitempty"`
	Time             primitive.DateTime    `json:"time,omitempty"`
	Status           string                `json:"status,omitempty"`
	PriceListRef     PriceList             `json:"pricelistref,omitempty"`
	ExpirationPeriod int64                 `json:"expirationperiod,omitempty"`
	Products         []InvoiceRowPopulated `json:"products,omitempty"`
	TotalProduct     float64               `json:"totalproduct,omitempty"`
	TotalTax         float64               `json:"totaltax,omitempty"`
	Total            float64               `json:"total,omitempty"`
	InventoryRef     Inventory             `json:"inventoryref,omitempty"` // for order
	Type             string                `json:"type,omitempty"`
	QuotationRef     Sales                 `json:"quotationref,omitempty"`
	SalesOrderRef    Sales                 `json:"salesorderref,omitempty"`
	ConvertedToSales bool                  `json:"convertedtosales,omitempty"`
	IsDelivered      bool                  `json:"isdelivered,omitempty"`
}

func (obj *SalesPopulated) CloneFrom(other Sales) {
	obj.ID = other.ID
	obj.Serial = other.Serial
	obj.CustomerRef = Contact{}
	obj.Date = other.Date
	obj.Time = other.Time
	obj.Status = other.Status
	obj.PriceListRef = PriceList{}
	obj.ExpirationPeriod = other.ExpirationPeriod
	obj.Products = []InvoiceRowPopulated{}
	obj.TotalProduct = other.TotalProduct
	obj.TotalTax = other.TotalTax
	obj.Total = other.Total
	obj.InventoryRef = Inventory{}
	obj.Type = other.Type
	obj.QuotationRef = Sales{}
	obj.SalesOrderRef = Sales{}
	obj.ConvertedToSales = other.ConvertedToSales
	obj.IsDelivered = other.IsDelivered
}

type SalesSearch struct {
	Type               string             `json:"type,omitempty"`
	TypeIsUsed         bool               `json:"typeisused,omitempty"`
	Status             string             `json:"status,omitempty"`
	StatusIsUsed       bool               `json:"statusisused,omitempty"`
	CustomerRef        primitive.ObjectID `json:"customerref,omitempty"`
	CustomerRefIsUsed  bool               `json:"customerrefisused,omitempty"`
	InventoryRef       primitive.ObjectID `json:"inventoryref,omitempty"`
	InventoryRefIsUsed bool               `json:"inventoryrefisused,omitempty"`
	DateFrom           primitive.DateTime `json:"datefrom,omitempty"`
	DateTo             primitive.DateTime `json:"dateto,omitempty"`
	DateRangeIsUsed    bool               `json:"daterangeisused,omitempty"`
	TimeFrom           primitive.DateTime `json:"timefrom,omitempty"`
	TimeTo             primitive.DateTime `json:"timeto,omitempty"`
	TimeRangeIsUsed    bool               `json:"timerangeisused,omitempty"`
	IsDelivered        bool               `json:"isdelivered,omitempty"`
	IsDeliveredIsUsed  bool               `json:"isdeliveredisused,omitempty"`
}

func (obj SalesSearch) GetSalesSearchBSONObj() bson.M {
	self := bson.M{}

	if obj.TypeIsUsed {
		self["type"] = obj.Type
	}

	if obj.StatusIsUsed {
		self["status"] = obj.Status
	}

	if obj.CustomerRefIsUsed {
		self["customerref"] = obj.CustomerRef
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

	if obj.IsDeliveredIsUsed {
		self["isdelivered"] = obj.IsDelivered
	}

	return self
}

func (obj Sales) GetModifcationBSONObj() bson.M {
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
