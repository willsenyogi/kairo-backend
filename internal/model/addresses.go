package model

import (
	"time"
)

type Address struct {
	ID        uint       `gorm:"primaryKey"`
	UserID    uint       `gorm:"not null;index"`
	User      User       `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL;"`
	Street    string     `gorm:"type:text;not null"`
	CityID    uint       `gorm:"index"`
	City      MasterCity `gorm:"foreignKey:CityID;references:ID;constraint:OnDelete:SET NULL;"`
	ZipCode   string     `gorm:"type:text;not null"`
	Label     string     `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
