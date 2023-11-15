package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Character struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Image  string             `json:"image_url" bson:"image_url"`
	Name   string             `json:"name" bson:"name"`
	Status string             `json:"status" bson:"status"`
	Birth  string             `json:"birth" bson:"birth"`
}
