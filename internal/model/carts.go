package model

import (
	"time"
)

type Cart struct {
	ID        uint        `gorm:"primaryKey"`
	UserID    uint        `gorm:"not null;index"`
	User      User        `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL;"`
	createdAt time.Time
	UpdatedAt time.Time
}