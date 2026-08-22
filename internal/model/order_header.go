package model

import (
	"time"

	"kairo-backend/internal/enum"
)

type OrderHeader struct {
	ID                uint        `gorm:"primaryKey"`
	UserID            uint        `gorm:"not null;index"`
	User              User        `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL;"`
	SellerID          uint        `gorm:"not null;index"`
	Seller            User        `gorm:"foreignKey:SellerID;references:ID;constraint:OnDelete:SET NULL;"`
	Status            enum.Status `gorm:"type:varchar(20);not null"`
	SubTotal          float64     `gorm:"type:decimal(10,2);not null"`
	ShippingCost      float64     `gorm:"type:decimal(10,2);not null"`
	Tax               float64     `gorm:"type:decimal(10,2);not null"`
	Total             float64     `gorm:"type:decimal(10,2);not null"`
	ShippingAddressID uint        `gorm:"not null;index"`
	ShippingAddress   Address     `gorm:"foreignKey:ShippingAddressID;references:ID;constraint:OnDelete:SET NULL;"`
	TransactionID     string      `gorm:"type:varchar(100);uniqueIndex;not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
