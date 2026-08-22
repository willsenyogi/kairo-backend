package model

import (
	"time"
)

type OrderDetail struct {
	ID         uint        `gorm:"primaryKey"`
	OrderID    uint        `gorm:"not null;index"`
	Order      OrderHeader `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE;"`
	ProductID  uint        `gorm:"not null;index"`
	Product    Product    `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE;"`
	unitPrice  float64     `gorm:"not null"`
	Quantity   uint        `gorm:"not null"`
	totalPrice float64     `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
