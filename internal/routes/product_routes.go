package routes

import (
	"kairo-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(r *gin.Engine, authHandler *handler.ProductHandler) {
	products := r.Group("/product")

	{
		products.GET("/", authHandler.GetAllProducts)
	}
}
