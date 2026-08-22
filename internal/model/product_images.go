package model

import (
	"time"
)

type ProductImage struct {
	ID        uint        `gorm:"primaryKey"`
	ProductID uint        `gorm:"not null;index"`
	Products  Product     `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE;"`
	URL       string      `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}