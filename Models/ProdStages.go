package Models

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProdStages struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Status      bool               `json:"status,omitempty"`
	Machine     bool               `json:"machine,omitempty"`
	MachineName string             `json:"machinename,omitempty"`
	Human       bool               `json:"human,omitempty"`
	HumanCount  int64              `json:"humancount,omitempty"`
	CostHours   float64            `json:"costhours,omitempty"`
}

func (obj ProdStages) Validate() error {
	return validation.ValidateStruct(&obj,
		validation.Field(&obj.Name, validation.Required),
	)
}

func (obj ProdStages) GetBSONModificationObj() bson.M {
	self := bson.M{
		"name":        obj.Name,
		"status":      obj.Status,
		"machine":     obj.Machine,
		"machinename": obj.MachineName,
		"human":       obj.Human,
		"humancount":  obj.HumanCount,
		"costhours":   obj.CostHours,
	}

	return self
}

type ProdStageSearch struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	IDIsUsed          bool               `json:"idisused,omitempty" bson:"idisused,omitempty"`
	Name              string             `json:"name,omitempty"`
	NameIsUsed        bool               `json:"nameisused,omitempty"`
	Status            bool               `json:"status,omitempty"`
	StatusIsUsed      bool               `json:"statusisused,omitempty"`
	Machine           bool               `json:"machine,omitempty"`
	MachineIsUsed     bool               `json:"machineisused,omitempty"`
	MachineName       string             `json:"machinename,omitempty"`
	MachineNameIsUsed bool               `json:"machinenameisused,omitempty"`
	Human             bool               `json:"human,omitempty"`
	HumanIsUsed       bool               `json:"humanisused,omitempty"`
	HumanCount        int64              `json:"humancount,omitempty"`
	HumanCountIsUsed  bool               `json:"humancountisused,omitempty"`
	CostHours         float64            `json:"costhours,omitempty"`
	CostHoursIsUsed   bool               `json:"costhoursisused,omitempty"`
}

func (obj ProdStageSearch) GetProdStagesSearchBSONObj() bson.M {
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
	if obj.MachineIsUsed {
		self["machine"] = obj.Machine
	}
	if obj.MachineNameIsUsed {
		regexPattern := fmt.Sprintf(".*%s.*", obj.MachineName)
		self["machinename"] = bson.D{{"$regex", primitive.Regex{Pattern: regexPattern, Options: "i"}}}
	}
	if obj.HumanIsUsed {
		self["human"] = obj.Human
	}
	if obj.HumanCountIsUsed {
		self["humancount"] = obj.HumanCount
	}

	if obj.CostHoursIsUsed {
		self["costhours"] = obj.CostHours
	}

	return self
}
