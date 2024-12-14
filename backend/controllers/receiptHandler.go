package controllers

import (
	"github.com/OAthooh/BiasharaTrack.git/models"
	"github.com/OAthooh/BiasharaTrack.git/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReceiptHandler struct {
	db *gorm.DB
}

func NewReceiptHandler(db *gorm.DB) *ReceiptHandler {
	return &ReceiptHandler{db: db}
}

// GetReceipt retrieves a receipt by its receipt number for a specific user
func (rh *ReceiptHandler) GetReceipt(c *gin.Context) {
	userID := c.GetUint("userID")
	receiptNumber := c.Param("receiptNumber")
	if receiptNumber == "" {
		utils.ErrorLogger("Receipt number is required")
		c.JSON(400, gin.H{"error": "Receipt number is required"})
		return
	}

	var receipt models.Receipt
	if err := rh.db.Preload("Items").Where("receipt_number = ? AND user_id = ?", receiptNumber, userID).First(&receipt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorLogger("Receipt not found for user: %d, receipt: %s", userID, receiptNumber)
			c.JSON(404, gin.H{"error": "Receipt not found"})
			return
		}
		utils.ErrorLogger("Failed to fetch receipt for user: %d, error: %v", userID, err)
		c.JSON(500, gin.H{"error": "Failed to fetch receipt"})
		return
	}

	utils.InfoLogger("Successfully fetched receipt for user: %d, receipt: %s", userID, receiptNumber)
	c.JSON(200, receipt)
}

// GetAllReceipts retrieves all receipts for a specific user with optional date range filtering
func (rh *ReceiptHandler) GetAllReceipts(c *gin.Context) {
	userID := c.GetUint("userID")
	var receipts []models.Receipt
	
	query := rh.db.Preload("Items").Where("user_id = ?", userID).Order("created_at desc")

	// Get date range filters from query params if they exist
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	if startDate != "" && endDate != "" {
		query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Find(&receipts).Error; err != nil {
		utils.ErrorLogger("Failed to fetch receipts for user: %d, error: %v", userID, err)
		c.JSON(500, gin.H{"error": "Failed to fetch receipts"})
		return
	}

	utils.InfoLogger("Successfully fetched all receipts for user: %d", userID)
	c.JSON(200, receipts)
}

// UpdateReceipt updates an existing receipt for a specific user
func (rh *ReceiptHandler) UpdateReceipt(c *gin.Context) {
	userID := c.GetUint("userID")
	receiptNumber := c.Param("receiptNumber")
	if receiptNumber == "" {
		utils.ErrorLogger("Receipt number is required")
		c.JSON(400, gin.H{"error": "Receipt number is required"})
		return
	}

	var receipt models.Receipt
	if err := rh.db.Where("receipt_number = ? AND user_id = ?", receiptNumber, userID).First(&receipt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorLogger("Receipt not found for user: %d, receipt: %s", userID, receiptNumber)
			c.JSON(404, gin.H{"error": "Receipt not found"})
			return
		}
		utils.ErrorLogger("Failed to fetch receipt for user: %d, error: %v", userID, err)
		c.JSON(500, gin.H{"error": "Failed to fetch receipt"})
		return
	}

	// Parse update data from request body
	var updateData models.Receipt
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorLogger("Invalid request body for user: %d, error: %v", userID, err)
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Update receipt fields
	if err := rh.db.Model(&receipt).Updates(updateData).Error; err != nil {
		utils.ErrorLogger("Failed to update receipt for user: %d, error: %v", userID, err)
		c.JSON(500, gin.H{"error": "Failed to update receipt"})
		return
	}

	utils.InfoLogger("Successfully updated receipt for user: %d, receipt: %s", userID, receiptNumber)
	c.JSON(200, receipt)
}
