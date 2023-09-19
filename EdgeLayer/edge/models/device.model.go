package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
{ version: 111,
  id: 1459633408,
  timestamp: 1685407036,
  input1: 1,
  input2: 1,
  output: 0,
  alarm_battery: false,
  alarm_power: false,
  sensor_error: false,
  sensors: [ { type: 0, value: 24.870001 }, { type: 1, value: 65.970001 } ] }
*/

type Sensors struct {
	Type  int32   `bson:"type" json:"type" `
	Value float32 `bson:"value" json:"value"`
}

type Device struct {
	InsertedID    string    `json:"InsertedID"`
	IdDevice      int32     `json:"id" validate:"required"`
	Input1        int32     `json:"input1, Input1"`
	Input2        int32     `json:"input2, Input2"`
	Output        int32     `json:"output, Output"`
	Alarm_battery bool      `json:"alarm_battery, Alarm_battery"`
	Alarm_power   bool      `json:"alarm_power, Alarm_power" `
	Sensor_error  bool      `json:"sensor_error, Sensor_error" `
	Sensors_obj   []Sensors `bson:"sensors,Sensors" json:"sensors, Sensors"`
	LastUpdated   time.Time `bson:"lastUpdated" json:"lastUpdated"`
}

func (note *Device) BeforeCreate(tx *gorm.DB) (err error) {
	note.InsertedID = uuid.New().String()
	return nil
}

func (note *Devices) BeforeCreate(tx *gorm.DB) (err error) {
	note.InsertedID = uuid.New().String()
	return nil
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

type Devices struct {
	InsertedID    string  `json:"InsertedID"`
	IdEdge        int32   `json:"idEdge" validate:"required"`
	IdNodo        int32   `json:"IdNodo" validate:"required"`
	Input1        int32   `json:"input1, Input1"`
	Input2        int32   `json:"input2, Input2"`
	Output        int32   `json:"output, Output"`
	Alarm_battery bool    `json:"alarm_battery, Alarm_battery"`
	Alarm_power   bool    `json:"alarm_power, Alarm_power" `
	Sensor_error  bool    `json:"sensor_error, Sensor_error" `
	Temperatura   float32 `bson:"temperatura,Temperatura" json:"temperatura,Temperatura"`
	Umidade       float32 `bson:"umidade,Umidade" json:"umidade,Umidade"`
	CreatedAt     string  `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt     string  `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

type DeviceHash struct {
	Id     int32  `json:"id" validate:"required"`
	IdEdge int32  `json:"idEdge"`
	IdNodo int32  `json:"IdNodo"`
	Hash   string `json:"hash"`
	Date   string `json:"date"`
}

type DeviceEdge struct {
	IdEdge int32  `json:"idEdge" validate:"required"`
	Ip     string `json:"ip"`
	User   string `json:"user"`
	Date   string `json:"date"`
	Hash   string `json:"hash"`
}

type DeviceNodo struct {
	IdNodo int32  `json:"id" validate:"required"`
	IdEdge int32  `json:"idEdge" validate:"required"`
	Period int32  `json:"period"`
	Date   string `json:"date"`
	Hash   string `json:"hash"`
}

type AlertRequest struct {
	Id          int32  `json:"id" validate:"required"`
	IdEdge      int32  `json:"idEdge"`
	IdNodo      int32  `json:"idNodo"`
	Hash        string `json:"hash"`
	Type        int32  `json:"type"`
	Date        string `json:"date"`
	Temperatura string `json:"temperatura"`
	Umidade     string `json:"umidade"`
	Rele        string `json:"rele"`
	Description string `json:"description"`
}
