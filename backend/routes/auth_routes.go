package routes

import (
	"polling-app/backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	authRoutes := r.Group("/api/auth")

	authRoutes.POST("/signup", authHandler.Signup)
	authRoutes.POST("/login", authHandler.Login)
}