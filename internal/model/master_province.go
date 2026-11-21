package model

import "time"

type MasterProvince struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"type:varchar(2);not null;uniqueIndex"`
	Name string `gorm:"type:varchar(100);not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
