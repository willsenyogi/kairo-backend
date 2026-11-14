package model

import (
	"time"

	"kairo-backend/internal/enum"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string    `gorm:"type:text;not null" json:"-"`
	Role      enum.Role `gorm:"type:integer;default:1"`
	AvatarURL string    `gorm:"type:text"`
	Seller    *Seller   `gorm:"foreignKey:UserID;references:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
