package pubmodel

import (
	"time"

	"gorm.io/gorm"
)

type PubPackageModel struct {
	Name      string            `json:"name" gorm:"not null;primaryKey;"`
	Private   bool              `json:"private" gorm:"not null;default:true"`
	Versions  []PubVersionModel `json:"versions" gorm:"foreignKey:PackageName;references:Name"`
	CreatedAt *time.Time        `json:"created_at,omitempty" gorm:"not null;"`
	UpdatedAt *time.Time        `json:"updated_at,omitempty" gorm:"not null;"`
	DeletedAt *gorm.DeletedAt   `json:"deleted_at,omitempty" gorm:"index"`
}

func (PubPackageModel) TableName() string {
	return "pub_packages"
}
