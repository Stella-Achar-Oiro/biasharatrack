package routes

import (
	"github.com/OAthooh/BiasharaTrack.git/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupReceiptRoutes(router *gin.Engine, db *gorm.DB) {
	receiptHandler := controllers.NewReceiptHandler(db)
	router.GET("/get-all-receipts", receiptHandler.GetAllReceipts)
}
