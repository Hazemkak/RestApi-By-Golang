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

type Contact struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty"`
	Status     bool               `json:"status,omitempty"`
	CompanyRef primitive.ObjectID `json:"companyref,omitempty" bson:"companyref,omitempty"`
	Phone      string             `json:"phone,omitempty"`
	Mobile     string             `json:"mobile,omitempty"`
	Email      string             `json:"email,omitempty"`
	Website    string             `json:"website,omitempty"`
	Address    string             `json:"address,omitempty"`
	JobTitle   string             `json:"jobtitle,omitempty"`
	Title      string             `json:"title,omitempty"`
	IsCompany  bool               `json:"iscompany,omitempty"`
	IsCustomer bool               `json:"iscustomer,omitempty"`
}

func (obj Contact) GetIdString() string {
	return obj.ID.String()
}

func (obj Contact) GetId() primitive.ObjectID {
	return obj.ID
}

func (obj Contact) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
	)
}
func (obj Contact) GetModifcationBSONObj() bson.M {
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

type ContactSearch struct {
	IDIsUsed         bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	NameIsUsed       bool               `json:"nameisused,omitempty"`
	Name             string             `json:"name,omitempty"`
	StatusIsUsed     bool               `json:"statusisused,omitempty"`
	Status           bool               `json:"status,omitempty"`
	CompanyRefIsUsed bool               `json:"companyrefisused,omitempty" bson:"companyrefisused,omitempty"`
	CompanyRef       primitive.ObjectID `json:"companyref,omitempty" bson:"companyref,omitempty"`
	PhoneIsUsed      bool               `json:"phoneisused,omitempty"`
	Phone            string             `json:"phone,omitempty"`
	MobileIsUsed     bool               `json:"mobileisused,omitempty"`
	Mobile           string             `json:"mobile,omitempty"`
	EmailIsUsed      bool               `json:"emailisused,omitempty"`
	Email            string             `json:"email,omitempty"`
	WebsiteIsUsed    bool               `json:"websiteisused,omitempty"`
	Website          string             `json:"website,omitempty"`
	AddressIsUsed    bool               `json:"addressisused,omitempty"`
	Address          string             `json:"address,omitempty"`
	JobTitleIsUsed   bool               `json:"jobtitleisused,omitempty"`
	JobTitle         string             `json:"jobtitle,omitempty"`
	TitleIsUsed      bool               `json:"titleisused,omitempty"`
	Title            string             `json:"title,omitempty"`
	IsCompanyIsUsed  bool               `json:"iscompanyisused,omitempty"`
	IsCompany        bool               `json:"iscompany,omitempty"`
	IsCustomerIsUsed bool               `json:"iscustomerisused,omitempty"`
	IsCustomer       bool               `json:"iscustomer,omitempty"`
}

func (obj ContactSearch) GetContactSearchBSONObj() bson.M {
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

	if obj.CompanyRefIsUsed {
		self["companyref"] = obj.CompanyRef
	}

	if obj.PhoneIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Phone)
		self["phone"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.MobileIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Mobile)
		self["mobile"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.EmailIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Email)
		self["email"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.WebsiteIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Website)
		self["website"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.AddressIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Address)
		self["address"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.JobTitleIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.JobTitle)
		self["jobtitle"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.TitleIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.Title)
		self["title"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}

	if obj.IsCompanyIsUsed {
		self["iscompany"] = obj.IsCompany
	}

	if obj.IsCustomerIsUsed {
		self["iscustomer"] = obj.IsCustomer
	}

	return self
}

type ContactPopulated struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name,omitempty"`
	Status     bool               `json:"status,omitempty"`
	CompanyRef Contact            `json:"companyref,omitempty" bson:"companyref,omitempty"`
	Phone      string             `json:"phone,omitempty"`
	Mobile     string             `json:"mobile,omitempty"`
	Email      string             `json:"email,omitempty"`
	Website    string             `json:"website,omitempty"`
	Address    string             `json:"address,omitempty"`
	JobTitle   string             `json:"jobtitle,omitempty"`
	Title      string             `json:"title,omitempty"`
	IsCompany  bool               `json:"iscompany,omitempty"`
	IsCustomer bool               `json:"iscustomer,omitempty"`
}

func (obj *ContactPopulated) CloneFrom(other Contact) {
	obj.ID = other.ID
	obj.Name = other.Name
	obj.Status = other.Status
	obj.CompanyRef = Contact{}
	obj.Phone = other.Phone
	obj.Mobile = other.Mobile
	obj.Email = other.Email
	obj.Website = other.Website
	obj.Address = other.Address
	obj.JobTitle = other.JobTitle
	obj.Title = other.Title
	obj.IsCompany = other.IsCompany
	obj.IsCustomer = other.IsCustomer
}
