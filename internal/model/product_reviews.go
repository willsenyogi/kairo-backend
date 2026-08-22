package model

import (
	"time"
)

type ProductReview struct {
	ID        uint     `gorm:"primaryKey"`
	ProductID uint     `gorm:"not null;index"`
	Products  Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE;"`
	UserID    uint     `gorm:"not null;index"`
	User      User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL;"`
	Rating    uint     `gorm:"not null"`
	Comment   string   `gorm:"type:text"`
	CreatedAt time.Time
}
