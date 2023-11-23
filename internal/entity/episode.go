package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Episode struct {
	ID        primitive.ObjectID
	Image     string
	Name      string
	Airing    string
	Writers   []string
	Directors string
	Season    int32
	Episode   int32
}
