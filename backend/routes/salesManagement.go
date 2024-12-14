package routes

import (
	"github.com/OAthooh/BiasharaTrack.git/controllers"
	"github.com/OAthooh/BiasharaTrack.git/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SalesManagementRoutes(router *gin.Engine, db *gorm.DB) {
	sm := controllers.NewSalesManagementHandler(db)

	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.POST("/record-sale", sm.SellProducts)
		authenticated.GET("/sales-history", sm.FetchSalesHistory)
		authenticated.GET("/sales-metrics", sm.FetchSalesMetrics)
	}
}
