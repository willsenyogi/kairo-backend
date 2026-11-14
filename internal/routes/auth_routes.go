package routes

import (
	"kairo-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	internalAdmin := r.Group("internal/auth")
	auth := r.Group("/auth/users")

	{
		internalAdmin.GET("/fetch-all", authHandler.GetAllUsers)
	}

	{
		auth.POST("new-user/register", authHandler.RegisterNewUser)
		auth.POST("login/validate-creds", authHandler.LoginValidateCreds)
	}

}
