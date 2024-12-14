package routes

import (
	"github.com/OAthooh/BiasharaTrack.git/controllers"
	"github.com/OAthooh/BiasharaTrack.git/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupReceiptRoutes(router *gin.Engine, db *gorm.DB) {
	receiptHandler := controllers.NewReceiptHandler(db)
	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.GET("/get-all-receipts", receiptHandler.GetAllReceipts)
	}
}
