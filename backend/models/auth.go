package models

import "time"

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Register struct {
	FullName     string `json:"fullName"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	BusinessName string `json:"businessName"`
	Telephone    string `json:"telephone"`
	Location     string `json:"location"`
}
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FullName     string    `gorm:"not null" json:"fullName"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password     string    `gorm:"not null" json:"password"`
	BusinessName string    `gorm:"not null" json:"businessName"`
	Telephone    string    `gorm:"not null" json:"telephone"`
	Location     string    `gorm:"not null" json:"location"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
