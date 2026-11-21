package model

import "time"

type MasterDistrict struct {
	ID     uint   `gorm:"primaryKey"`
	CityID uint   `gorm:"not null;index"`
	Code   string `gorm:"type:varchar(7);not null;uniqueIndex"`
	Name   string `gorm:"type:varchar(100);not null"`

	City MasterCity `gorm:"foreignKey:CityID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
