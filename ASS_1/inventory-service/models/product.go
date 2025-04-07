package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Category string             `json:"category" bson:"category"`
	Stock    int                `json:"stock" bson:"stock"`
	Price    float64            `json:"price" bson:"price"`
}
