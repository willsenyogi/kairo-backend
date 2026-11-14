
package model

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TableName string         `gorm:"type:varchar(100);not null" json:"table_name"`
	RecordID  string         `gorm:"type:varchar(100);not null" json:"record_id"`
	Action    string         `gorm:"type:varchar(20);not null" json:"action"`
	OldValues datatypes.JSON `gorm:"type:jsonb" json:"old_values,omitempty"`
	NewValues datatypes.JSON `gorm:"type:jsonb" json:"new_values,omitempty"`
	ActorID   *uint          `gorm:"index" json:"actor_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
}
