package routes

import (
	"github.com/OAthooh/BiasharaTrack.git/controllers"
	"github.com/OAthooh/BiasharaTrack.git/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreditRoutes(router *gin.Engine, db *gorm.DB) {
	cm := controllers.NewCreditManager(db)

	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.GET("/credit-history", cm.GetCreditsHistory)
	}
}
