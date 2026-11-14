package main

import (
	"log"

	"kairo-backend/internal/config"
	"kairo-backend/internal/database"
	"kairo-backend/internal/handler"
	"kairo-backend/internal/model"
	"kairo-backend/internal/repository"
	"kairo-backend/internal/routes"
	"kairo-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewDB(cfg)

	AUTOMIGRATE := true

	if AUTOMIGRATE {
		if err := db.AutoMigrate(
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
			&model.AuditLog{},
		); err != nil {
			log.Fatal("Auto migration failed: ", err)
		}
	}

	r := gin.Default()

	// auth service + handler
	authRepo := repository.InitAuthRepository(db)
	authService := service.InitAuthService(authRepo)
	authHandler := handler.InitAuthHandler(authService)

	// wire product service + handler
	productRepo := repository.InitProductRepository(db)
	productService := service.InitProductService(productRepo)
	productHandler := handler.InitProductHandler(productService)

	// routes initialization
	routes.AuthRoutes(r, authHandler)
	routes.ProductRoutes(r, productHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
