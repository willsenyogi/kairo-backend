package main

import (
	"kairo-backend/internal/config"
	"kairo-backend/internal/database"
	"kairo-backend/internal/model"

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
			&model.Products{},
		)
	}

	
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API jalan",
		})
	})

	r.Run(":8080")
}
