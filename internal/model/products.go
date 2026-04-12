package model

import (
	"time"
)

type Products struct {
	ID        uint        `gorm:"primaryKey"`
	Name      string      `gorm:"type:varchar(100);not null"`
	SellerID  uint        `gorm:"not null"`
	Seller    Seller      `gorm:"foreignKey:SellerID;references:ID;constraint:OnDelete:CASCADE;"`
	CategoryID uint       `gorm:"not null"`
	Category  Category    `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE;"`
	Description string    `gorm:"type:text"`
	Price	   float64   `gorm:"not null"`
	DiscountedPrice float64   `gorm:"default:null"`
	Stock	   uint      `gorm:"not null"`
	SoldCount uint      `gorm:"default:0"`
	Rating    float64   `gorm:"default:0"`
	InStock   bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}