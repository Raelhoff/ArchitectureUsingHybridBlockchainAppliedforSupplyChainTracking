package models

import "go.mongodb.org/mongo-driver/bson/primitive"
//import "time"

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

type Devices struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
    	IdEdge        int32   `json:"idedge" validate:"required"`
	IdNodo        int32   `json:"idnodo" validate:"required"`
	Input1        int32   `json:"input1"`
	Input2        int32   `json:"input2"`
	Output        int32   `json:"output"`
	Hash          string  `json:"hash"`
	Alarm_battery bool    `json:"alarm_battery, Alarm_battery"`
	Alarm_power   bool    `json:"alarm_power, Alarm_power" `
	Sensor_error  bool    `json:"sensor_error, Sensor_error" `
	Temperatura   float32 `bson:"temperatura,Temperatura" json:"temperatura,Temperatura"`
	Umidade       float32 `bson:"umidade,Umidade" json:"umidade,Umidade"`
	CreatedAt     string  `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	UpdatedAt     string  `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}
