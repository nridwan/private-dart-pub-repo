package base

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;not null;primaryKey;default:uuid_generate_v4()"`
	CreatedAt *time.Time      `json:"created_at,omitempty" gorm:"not null;"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" gorm:"not null;"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
