package main

import (
	"github.com/gin-gonic/gin"
	"order-service/db"
	"order-service/internal/order"
)

func main() {
	client := db.InitDB()

	repo := &order.Repository{Client: client}

	service := &order.Service{Repo: repo}

	handler := &order.Handler{Service: service}

	r := gin.Default()

	handler.RegisterRoutes(r)

	r.Run(":8082")
}
