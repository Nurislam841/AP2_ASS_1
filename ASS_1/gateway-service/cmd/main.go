package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Order struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID string             `json:"user_id" bson:"user_id"`
	Status string             `json:"status" bson:"status"`
}

var orderCollection *mongo.Collection

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	orderCollection = client.Database("orderdb").Collection("orders")

	router := gin.Default()

	router.POST("/orders", createOrder)
	router.GET("/orders", getOrdersByUser)
	router.GET("/orders/:id", getOrderByID)
	router.PATCH("/orders/:id", updateOrderStatus)

	router.Run(":8080")
}

func createOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := orderCollection.InsertOne(context.TODO(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	order.ID = res.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, order)
}

func getOrdersByUser(c *gin.Context) {
	userID := c.Query("user_id")
	filter := bson.M{}
	if userID != "" {
		filter["user_id"] = userID
	}

	cursor, err := orderCollection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO()) // Закрытие курсора

	var orders []Order
	if err := cursor.All(context.TODO(), &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func getOrderByID(c *gin.Context) {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var order Order
	err = orderCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func updateOrderStatus(c *gin.Context) {
	idParam := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{"$set": bson.M{"status": body.Status}}
	_, err = orderCollection.UpdateByID(context.TODO(), objectID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}
