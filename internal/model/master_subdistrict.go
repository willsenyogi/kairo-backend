package model

import "time"

type MasterSubdistrict struct {
	ID         uint   `gorm:"primaryKey"`
	DistrictID uint   `gorm:"not null;index"`
	Code       string `gorm:"type:varchar(10);not null;uniqueIndex"`
	Name       string `gorm:"type:varchar(100);not null"`

	District MasterDistrict `gorm:"foreignKey:DistrictID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
