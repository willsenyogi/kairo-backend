package model

import "time"

type MasterCity struct {
	ID         uint   `gorm:"primaryKey"`
	ProvinceID uint   `gorm:"not null;index"`
	Code       string `gorm:"type:varchar(4);not null;uniqueIndex"`
	Name       string `gorm:"type:varchar(100);not null"`

	Province MasterProvince `gorm:"foreignKey:ProvinceID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
