package routes

import (
	"github.com/OAthooh/BiasharaTrack.git/controllers"
	"github.com/OAthooh/BiasharaTrack.git/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InventoryManagementRoutes(router *gin.Engine, db *gorm.DB) {
	im := controllers.NewInventoryManagementHandler(db)

	// Protected routes
	authenticated := router.Group("/")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.POST("/create-product", im.CreateProduct)
		authenticated.PUT("/update-product/:id", im.UpdateProduct)
		authenticated.DELETE("/delete-product/:id", im.DeleteProduct)
		authenticated.GET("/get-product/:id", im.GetProduct)
		authenticated.GET("/get-all-products", im.GetAllProducts)
		authenticated.GET("/get-low-stock-alerts", im.GetLowStockAlerts)
		authenticated.GET("/lookup-barcode/:barcode", im.LookupBarcode)
		authenticated.GET("/search-products", im.SearchProducts)
	}

	// Public routes (if any)
	router.Static("/uploads", "./uploads")
}
