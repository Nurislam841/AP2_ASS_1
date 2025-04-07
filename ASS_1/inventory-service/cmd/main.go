package main

import (
	"github.com/gin-gonic/gin"
	"inventory-service/db"
	"inventory-service/internal/product"
)

func main() {
	database := db.InitDB()

	repo := &product.Repository{DB: database}

	service := &product.Service{Repo: repo}

	handler := &product.Handler{Service: service}

	r := gin.Default()

	handler.RegisterRoutes(r)

	r.Run(":8081")
}
