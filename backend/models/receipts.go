package models


import "time"

type Receipt struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ReceiptNumber string    `json:"receiptNumber" gorm:"unique;not null"`
	CustomerName  string    `json:"customerName"`
	Date          time.Time `json:"date"`
	PaymentMethod string    `json:"paymentMethod"`
	TotalAmount   float64   `json:"totalAmount"`
	Items         []Item    `json:"items" gorm:"foreignKey:ReceiptID"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	
}

type Item struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	ReceiptID  uint    `json:"receiptId"`
	ProductID  uint    `json:"productId"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unitPrice"`
	TotalPrice float64 `json:"totalPrice"`
}
