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

type Sensors struct {
	Type  int32   `json:"type" validate:"required"`
	Value float32 `json:"value" validate:"required"`
}

type Device struct {
	Id            primitive.ObjectID `json:"omitempty"`
	IdDevice      int32              `json:"id" validate:"required"`
	Input1        int32              `json:"input1" validate:"required"`
	Input2        int32              `json:"input2" validate:"required"`
	Output        int32              `json:"output" validate:"required"`
	Alarm_battery bool               `json:"alarm_battery" validate:"required"`
	Alarm_power   bool               `json:"alarm_power" validate:"required"`
	Sensor_error  bool               `json:"sensor_error" validate:"required"`
	Temp          float32            `json:"temp"`
	Umid          float32            `json:"umid"`
	LastUpdated   string             `bson:"lastUpdated,omitempty" json:"lastUpdated,omitempty"`
}

type Devices struct {
	Version   float32  `json:"version"`
	Id        int32    `json:"id"`
	Timestamp int64    `json:"timestamp"`
	Input1    int32    `json:"input1"`
	Input2    int32    `json:"input2"`
	Output    int32    `json:"output"`
	Sensors   []Sensor `json:"sensors"`
}

type Sensor struct {
	Type  int     `json:"type"`
	Value float32 `json:"value"`
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
