package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//import "time"

type DevicesNodo struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	IdNodo int32              `json:"IdNodo" validate:"required"`
	IdEdge int32              `json:"IdEdge" validate:"required"`
	Period int                `json:"Period"`
	Date   string             `gorm:"not null;default:'1970-01-01 00:00:01'" json:"Date,omitempty"`
	Hash   string             `json:"Hash"`
}
