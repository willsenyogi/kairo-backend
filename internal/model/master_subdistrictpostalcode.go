package model

import "time"

type MasterSubdistrictPostalCode struct {
	ID            uint `gorm:"primaryKey"`
	SubdistrictID uint `gorm:"not null;index"`
	PostalCodeID  uint `gorm:"not null;index"`

	Subdistrict MasterSubdistrict `gorm:"foreignKey:SubdistrictID"`
	PostalCode  MasterPostalCode  `gorm:"foreignKey:PostalCodeID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
