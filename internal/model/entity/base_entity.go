package entity

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type BaseEntity struct {
	CreatedAt *time.Time            `gorm:"column:created_at;autoCreateTime"`
	CreatedBy string                `gorm:"column:created_by;type:varchar(50)"`
	UpdatedAt *time.Time            `gorm:"column:updated_at;autoUpdateTime"`
	UpdatedBy string                `gorm:"column:updated_by;type:varchar(50)"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy string                `gorm:"column:deleted_by;type:varchar(50)"`
}
