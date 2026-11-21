package main

import (
	"log"

	"kairo-backend/configs"
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
	appConfig := configs.LoadAppConfig()
	dbConfig := config.LoadConfig()

	db := database.NewDB(dbConfig)

	log.Println("[INFO] Current running environment :", appConfig.AppEnv)

	if appConfig.AutoMigrate {
		log.Println("[WARNING] AutoMigrate is ENABLED. Database schema may be modified.")
		if err := db.AutoMigrate(model.AllModels()...); err != nil {
			log.Fatal("Auto migration failed: ", err)
		}
		log.Println("[INFO] AutoMigrate completed successfully.")
	} else {
		log.Println("[INFO] AutoMigrate is DISABLED. Skipping database migration.")
	}

	r := gin.Default()

	// Auth service + handler
	authRepo := repository.InitAuthRepository(db)
	authService := service.InitAuthService(authRepo)
	authHandler := handler.InitAuthHandler(authService)

	// Product service + handler
	productRepo := repository.InitProductRepository(db)
	productService := service.InitProductService(productRepo)
	productHandler := handler.InitProductHandler(productService)

	// Routes initialization
	routes.AuthRoutes(r, authHandler)
	routes.ProductRoutes(r, productHandler)

	if err := r.Run(":" + appConfig.ServerPort); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
