package controllers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/OAthooh/BiasharaTrack.git/models"
	"github.com/OAthooh/BiasharaTrack.git/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SalesManagementHandler struct {
	db *gorm.DB
}

func NewSalesManagementHandler(db *gorm.DB) *SalesManagementHandler {
	return &SalesManagementHandler{db: db}
}

// Define the structure for a single sell request
type SellRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required"`
	Note      string  `json:"note"`
	Amount    float64 `json:"amount"`
}

// Define the structure for the sale data from the front end
type SaleData struct {
	Products         []SellRequest `json:"products" binding:"required"`
	PaymentMethod    string        `json:"payment_method"`
	CustomerName     string        `json:"customer_name"`
	CustomerPhone    string        `json:"customer_phone"`
	ReferenceNumber  string        `json:"reference_number"`
	AmountPaid       float64       `json:"amount_paid"`
	RemainingBalance float64       `json:"remaining_balance"`
}

func generateReceiptNumber() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("RCP-%d%d", time.Now().Unix(), rand.Intn(1000))
}

func (im *SalesManagementHandler) SellProducts(c *gin.Context) {
	// Check if request body is empty
	if c.Request.Body == nil {
		utils.ErrorLogger("Request body is empty")
		c.JSON(400, gin.H{"error": "Request body is required"})
		return
	}

	// Parse the request body as SaleData
	var saleData SaleData
	if err := c.ShouldBindJSON(&saleData); err != nil {
		utils.ErrorLogger("Failed to parse request body: %v", err)
		c.JSON(400, gin.H{
			"error":   "Invalid request format. Expected JSON with sale data",
			"details": err.Error(),
		})
		return
	}

	// Validate the products array is not empty
	if len(saleData.Products) == 0 {
		utils.ErrorLogger("Empty products array received")
		c.JSON(400, gin.H{"error": "At least one product is required"})
		return
	}

	// Validate payment method and sale type
	if saleData.PaymentMethod == "" {
		utils.ErrorLogger("Incomplete sale request in payment method or sale type")
		c.JSON(400, gin.H{"error": "Payment method is required"})
		return
	}

	// Validate credit sale requirements
	if strings.ToUpper(saleData.PaymentMethod) == "CREDIT" {
		if saleData.CustomerName == "" || saleData.CustomerPhone == "" {
			utils.WarningLogger("Incomplete credit sale request for product %d", saleData.Products[0].ProductID)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Name and phone number are required for credit sales for product %d", saleData.Products[0].ProductID)})
			return
		}
	}

	// Get user ID from context (assuming it's set during authentication)
	userID := c.GetUint("userID")
	if userID == 0 {
		utils.ErrorLogger("User not authenticated")
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	utils.InfoLogger("Processing sale for user %d with payment method %s", userID, saleData.PaymentMethod)

	processSales(saleData, userID, im, c)
}

func processSales(saleData SaleData, userID uint, im *SalesManagementHandler, c *gin.Context) {
	// Start transaction
	tx := im.db.Begin()
	if tx.Error != nil {
		utils.ErrorLogger("Failed to start transaction: %v", tx.Error)
		c.JSON(500, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create receipt
	receipt := models.Receipt{
		UserID:        userID,
		ReceiptNumber: generateReceiptNumber(),
		CustomerName:  saleData.CustomerName,
		Date:          time.Now(),
		PaymentMethod: saleData.PaymentMethod,
		TotalAmount:   0, // Will be updated as we process items
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := tx.Create(&receipt).Error; err != nil {
		tx.Rollback()
		utils.ErrorLogger("Failed to create receipt: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create receipt"})
		return
	}

	for _, sellRequest := range saleData.Products {
		// Get current inventory and product
		var inventory models.Inventory
		var product models.Product
		if err := tx.Where("product_id = ? AND user_id = ?", sellRequest.ProductID, userID).First(&inventory).Error; err != nil {
			tx.Rollback()
			utils.ErrorLogger("Product not found in inventory: product_id= %d %v", sellRequest.ProductID, err)
			c.JSON(404, gin.H{"error": fmt.Sprintf("Product %d not found in inventory", sellRequest.ProductID)})
			return
		}

		if err := tx.Where("id = ? AND user_id = ?", sellRequest.ProductID, userID).First(&product).Error; err != nil {
			tx.Rollback()
			utils.ErrorLogger("Product not found: product_id= %d %v", sellRequest.ProductID, err)
			c.JSON(404, gin.H{"error": fmt.Sprintf("Product %d not found", sellRequest.ProductID)})
			return
		}

		// Check if we have enough stock
		if inventory.Quantity < sellRequest.Quantity {
			tx.Rollback()
			utils.WarningLogger("Insufficient stock for product %d. Requested: %d, Available: %d",
				sellRequest.ProductID, sellRequest.Quantity, inventory.Quantity)
			c.JSON(400, gin.H{"error": fmt.Sprintf("Insufficient stock for product %d", sellRequest.ProductID)})
			return
		}

		// Create receipt item
		item := models.Item{
			ReceiptID:  receipt.ID,
			ProductID:  sellRequest.ProductID,
			Name:       product.Name,
			Quantity:   sellRequest.Quantity,
			UnitPrice:  sellRequest.Amount / float64(sellRequest.Quantity),
			TotalPrice: sellRequest.Amount,
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			utils.ErrorLogger("Failed to create receipt item: %v", err)
			c.JSON(500, gin.H{"error": "Failed to create receipt item"})
			return
		}

		// Update receipt total
		receipt.TotalAmount += item.TotalPrice

		// Update inventory
		inventory.Quantity -= sellRequest.Quantity
		inventory.LastUpdated = time.Now()
		if err := tx.Save(&inventory).Error; err != nil {
			tx.Rollback()
			utils.ErrorLogger("Failed to update inventory for product %d: %v", sellRequest.ProductID, err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to update inventory for product %d", sellRequest.ProductID)})
			return
		}

		// Record stock movement
		stockMovement := models.StockMovement{
			UserID:         userID,
			ProductID:      sellRequest.ProductID,
			ChangeType:     saleData.PaymentMethod,
			QuantityChange: -sellRequest.Quantity,
			Note:           sellRequest.Note,
			CreatedAt:      time.Now(),
		}

		if err := tx.Create(&stockMovement).Error; err != nil {
			tx.Rollback()
			utils.ErrorLogger("Failed to create stock movement for product %d: %v", sellRequest.ProductID, err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to record stock movement for product %d", sellRequest.ProductID)})
			return
		}

		// Handle credit sale
		if strings.ToUpper(saleData.PaymentMethod) == "CREDIT" {
			creditTx := models.CreditTransaction{
				UserID:       userID,
				ProductID:    sellRequest.ProductID,
				Name:         saleData.CustomerName,
				PhoneNumber:  saleData.CustomerPhone,
				Quantity:     sellRequest.Quantity,
				CreditAmount: sellRequest.Amount,
				BalanceDue:   saleData.RemainingBalance,
				Status:       "PENDING",
			}

			if err := tx.Create(&creditTx).Error; err != nil {
				tx.Rollback()
				utils.ErrorLogger("Failed to create credit transaction for product %d: %v", sellRequest.ProductID, err)
				c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to record credit transaction for product %d", sellRequest.ProductID)})
				return
			}
		}

		// Record sales transaction
		salesTransaction := models.SalesTransaction{
			UserID:          userID,
			ProductID:       sellRequest.ProductID,
			Quantity:        sellRequest.Quantity,
			TotalAmount:     sellRequest.Amount,
			PaymentMethod:   saleData.PaymentMethod,
			CustomerName:    saleData.CustomerName,
			CustomerPhone:   saleData.CustomerPhone,
			ReferenceNumber: saleData.ReferenceNumber,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		if err := tx.Create(&salesTransaction).Error; err != nil {
			tx.Rollback()
			utils.ErrorLogger("Failed to create sales transaction for product %d: %v", sellRequest.ProductID, err)
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to record sales transaction for product %d", sellRequest.ProductID)})
			return
		}

		// Check for low stock alert
		if inventory.Quantity <= inventory.LowStockThreshold {
			alert := models.LowStockAlert{
				ProductID:    sellRequest.ProductID,
				AlertMessage: fmt.Sprintf("Product stock is low. Current quantity: %d", inventory.Quantity),
				Resolved:     false,
				CreatedAt:    time.Now(),
			}

			if err := tx.Create(&alert).Error; err != nil {
				utils.ErrorLogger("Failed to create low stock alert for product %d: %v", sellRequest.ProductID, err)
			}
		}
	}

	// Update receipt with final total
	if err := tx.Save(&receipt).Error; err != nil {
		tx.Rollback()
		utils.ErrorLogger("Failed to update receipt total: %v", err)
		c.JSON(500, gin.H{"error": "Failed to update receipt total"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.ErrorLogger("Failed to commit transaction: %v", err)
		c.JSON(500, gin.H{"error": "Failed to complete sales"})
		return
	}

	utils.InfoLogger("Successfully processed sales")
	c.JSON(200, gin.H{
		"message": "Sales recorded successfully",
	})
}

// Fetch sales history
func (im *SalesManagementHandler) FetchSalesHistory(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		utils.ErrorLogger("User not authenticated")
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var salesTransactions []struct {
		models.SalesTransaction
		ProductName string `json:"product_name"`
	}

	if err := im.db.Table("sales_transactions").
		Select("sales_transactions.*, products.name as product_name").
		Joins("JOIN products ON sales_transactions.product_id = products.id").
		Where("sales_transactions.user_id = ?", userID).
		Find(&salesTransactions).Error; err != nil {
		utils.ErrorLogger("Failed to fetch sales history: %v", err)
		c.JSON(500, gin.H{"error": "Failed to fetch sales history"})
		return
	}

	utils.InfoLogger("Successfully fetched sales history")
	c.JSON(200, salesTransactions)
}

// GetSalesMetrics returns sales metrics including daily, weekly and monthly revenue
func (im *SalesManagementHandler) FetchSalesMetrics(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		utils.ErrorLogger("User not authenticated")
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Get daily revenue
	var dailyRevenue float64
	if err := im.db.Model(&models.SalesTransaction{}).
		Where("user_id = ? AND created_at >= ?", userID, startOfDay).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&dailyRevenue).Error; err != nil {
		utils.ErrorLogger("Failed to fetch daily revenue: %v", err)
		c.JSON(500, gin.H{"error": "Failed to fetch sales metrics"})
		return
	}

	// Get weekly revenue
	var weeklyRevenue float64
	if err := im.db.Model(&models.SalesTransaction{}).
		Where("user_id = ? AND created_at >= ?", userID, startOfWeek).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&weeklyRevenue).Error; err != nil {
		utils.ErrorLogger("Failed to fetch weekly revenue: %v", err)
		c.JSON(500, gin.H{"error": "Failed to fetch sales metrics"})
		return
	}

	// Get monthly revenue
	var monthlyRevenue float64
	if err := im.db.Model(&models.SalesTransaction{}).
		Where("user_id = ? AND created_at >= ?", userID, startOfMonth).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&monthlyRevenue).Error; err != nil {
		utils.ErrorLogger("Failed to fetch monthly revenue: %v", err)
		c.JSON(500, gin.H{"error": "Failed to fetch sales metrics"})
		return
	}

	// Get payment method breakdown for current month
	var paymentBreakdown struct {
		Cash   float64
		Mpesa  float64
		Credit float64
	}
	if err := im.db.Model(&models.SalesTransaction{}).
		Where("user_id = ? AND created_at >= ?", userID, startOfMonth).
		Select(`
			COALESCE(SUM(CASE WHEN payment_method = 'cash' THEN total_amount ELSE 0 END), 0) as cash,
			COALESCE(SUM(CASE WHEN payment_method = 'mpesa' THEN total_amount ELSE 0 END), 0) as mpesa,
			COALESCE(SUM(CASE WHEN payment_method = 'credit' THEN total_amount ELSE 0 END), 0) as credit
		`).
		Scan(&paymentBreakdown).Error; err != nil {
		utils.ErrorLogger("Failed to fetch payment breakdown: %v", err)
		c.JSON(500, gin.H{"error": "Failed to fetch sales metrics"})
		return
	}

	// Get top products for current month
	type TopProduct struct {
		ProductName string  `json:"product_name"`
		Quantity    int     `json:"quantity"`
		Revenue     float64 `json:"revenue"`
	}
	var topProducts []TopProduct
	if err := im.db.Table("sales_transactions").
		Select("products.name as product_name, SUM(sales_transactions.quantity) as quantity, SUM(sales_transactions.total_amount) as revenue").
		Joins("JOIN products ON sales_transactions.product_id = products.id").
		Where("sales_transactions.user_id = ? AND sales_transactions.created_at >= ?", userID, startOfMonth).
		Group("products.id, products.name").
		Order("revenue DESC").
		Limit(1).
		Scan(&topProducts).Error; err != nil {
		utils.ErrorLogger("Failed to fetch top products: %v", err)
		c.JSON(500, gin.H{"error": "Failed to fetch sales metrics"})
		return
	}

	utils.InfoLogger("Successfully fetched sales metrics")
	c.JSON(200, gin.H{
		"dailyRevenue":           dailyRevenue,
		"weeklyRevenue":          weeklyRevenue,
		"monthlyRevenue":         monthlyRevenue,
		"topProducts":            topProducts,
		"paymentMethodBreakdown": paymentBreakdown,
	})
}
