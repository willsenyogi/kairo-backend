package model

import "time"

type MasterPostalCode struct {
	ID         uint   `gorm:"primaryKey"`
	PostalCode string `gorm:"type:varchar(5);not null;uniqueIndex"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
