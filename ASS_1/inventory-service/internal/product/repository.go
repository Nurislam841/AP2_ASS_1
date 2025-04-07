package product

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"inventory-service/models"
)

type Repository struct {
	DB *mongo.Collection
}

func (r *Repository) Create(product *models.Product) error {
	_, err := r.DB.InsertOne(context.Background(), product)
	return err
}

func (r *Repository) GetAll() ([]models.Product, error) {
	var products []models.Product
	cursor, err := r.DB.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) Update(id primitive.ObjectID, updates map[string]interface{}) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}
	_, err := r.DB.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *Repository) Delete(id primitive.ObjectID) error {
	_, err := r.DB.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *Repository) GetByID(id primitive.ObjectID) (*models.Product, error) {
	var product models.Product
	err := r.DB.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	return &product, err
}
