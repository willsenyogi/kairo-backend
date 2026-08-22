package model

import (
	"time"
)

type Payment struct {
	ID            uint          `gorm:"primaryKey"`
	OrderID       uint          `gorm:"not null;index"`
	Order         OrderHeader   `gorm:"foreignKey:OrderID;references:ID;constraint:OnDelete:CASCADE;"`
	Amount        float64       `gorm:"not null"`
	MethodID      uint          `gorm:"not null;index"`
	Method        PaymentMethod `gorm:"foreignKey:MethodID;references:ID;constraint:OnDelete:SET NULL;"`
	Status        string        `gorm:"type:varchar(50);not null"`
	transactionID string        `gorm:"type:varchar(100);not null"`
	paymentDate   time.Time     `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
