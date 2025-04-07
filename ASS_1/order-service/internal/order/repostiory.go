package order

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"order-service/models"

	"strconv"
)

type Repository struct {
	Client *mongo.Client
}

func (r *Repository) Create(order *models.Order) error {
	collection := r.Client.Database("order_db").Collection("orders")
	_, err := collection.InsertOne(context.Background(), order)
	return err
}

func (r *Repository) GetAll() ([]models.Order, error) {
	collection := r.Client.Database("order_db").Collection("orders")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *Repository) GetByID(id uint) (*models.Order, error) {
	collection := r.Client.Database("order_db").Collection("orders")
	filter := bson.M{"_id": primitive.ObjectIDFromHex(strconv.Itoa(int(id)))}
	var order models.Order
	err := collection.FindOne(context.Background(), filter).Decode(&order)
	return &order, err
}

func (r *Repository) UpdateStatus(id uint, status string) error {
	collection := r.Client.Database("order_db").Collection("orders")
	filter := bson.M{"_id": primitive.ObjectIDFromHex(strconv.Itoa(int(id)))}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *Repository) GetByUserID(userID uint) ([]models.Order, error) {
	collection := r.Client.Database("order_db").Collection("orders")
	filter := bson.M{"user_id": userID}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
