package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//import "time"

type DevicesEdge struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	IdEdge int32              `json:"IdEdge" validate:"required"`
	Ip     string             `json:"Ip"`
	User   string             `json:"User"`
	Date   string             `gorm:"not null;default:'1970-01-01 00:00:01'" json:"Date,omitempty"`
	Hash   string             `json:"Hash"`
}
