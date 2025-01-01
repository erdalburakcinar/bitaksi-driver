package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Coordinates struct {
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
}

type DriverWithDistance struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Location Location           `bson:"location" json:"location"`
	Distance float64            `bson:"distance" json:"distance"` // Distance from the "near" point
}

type Location struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}
