package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     uint               `json:"user_id" bson:"user_id"`
	ProductID  uint               `json:"product_id" bson:"product_id"`
	Quantity   int                `json:"quantity" bson:"quantity"`
	TotalPrice float64            `json:"total_price" bson:"total_price"`
	Status     string             `json:"status" bson:"status"`
}
