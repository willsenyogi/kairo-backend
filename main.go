package main

import (
	"kairo-backend/internal/config"
	"kairo-backend/internal/database"
	"kairo-backend/internal/handler"
	"kairo-backend/internal/model"
	"kairo-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewDB(cfg)

	AUTOMIGRATE := true

	if AUTOMIGRATE == true {
		db.AutoMigrate(
			&model.MasterCity{},
			&model.User{},
			&model.Seller{},
			&model.Category{},
			&model.Product{},
			&model.ProductImage{},
			&model.Address{},
			&model.OrderHeader{},
			&model.OrderDetail{},
			&model.Payment{},
		)
	}

	r := gin.Default()

	// wire product service + handler
	productService := service.NewProductService(db)
	productHandler := handler.ProductHandler{ProductService: productService}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API jalan"})
	})

	r.GET("/products", productHandler.GetAllProducts)

	r.Run(":8080")
}
