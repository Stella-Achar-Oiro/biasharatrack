package routes

import (
	"github.com/OAthooh/BiasharaTrack.git/controllers"
	"github.com/OAthooh/BiasharaTrack.git/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthRoutes sets up authentication related routes for the application
func AuthRoutes(router *gin.Engine, db *gorm.DB) {
	auth := controllers.NewAuthHandler(db)

	// Public routes
	router.POST("/login", auth.Login)
	router.POST("/register", auth.Register)

	// Protected routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.GET("/verify-token", auth.VerifyToken)
	}
}
