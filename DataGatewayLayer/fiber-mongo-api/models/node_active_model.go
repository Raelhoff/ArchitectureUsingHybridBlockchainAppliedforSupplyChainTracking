package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//import "time"

type NodeActive struct {
	_id          primitive.ObjectID `bson:"_id,omitempty"`
	IdEdge       int32              `json:"idedge" validate:"required"`
	IdNodo       int32              `json:"idnodo" validate:"required"`
	IsActive     bool               `json:"isactive"`
	Updateat     string             `gorm:"not null;default:'1970-01-01 00:00:01'" json:"updateat,omitempty"`
	Notification bool               `json:"notification"`
	Name         string             `json:"name"`
}
