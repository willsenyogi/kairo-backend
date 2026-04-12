package model

import (
	"time"
)

type Seller struct {
	ID          uint       `gorm:"primaryKey"`
	UserID      uint       `gorm:"not null;uniqueIndex"`
	User        User       `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
	Description string     `gorm:"type:text"`
	ShipsFrom   uint       `gorm:"not null"`
	MasterCity  MasterCity `gorm:"foreignKey:ShipsFrom;references:ID"`
	ReviewsCount uint       `gorm:"default:0"`
	Rating       float64    `gorm:"default:0"`
	SalesCount   uint       `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
