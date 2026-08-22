package model

import (
	"time"
)

type CartDetail struct {
	ID        uint        `gorm:"primaryKey"`
	cartID   uint        `gorm:"not null;index"`
	Cart      Cart        `gorm:"foreignKey:cartID;references:ID;constraint:OnDelete:CASCADE;"`
	ProductID uint        `gorm:"not null;index"`
	Product   Product    `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE;"`
	ProductPrice float64     `gorm:"not null"`
	Quantity  uint        `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}