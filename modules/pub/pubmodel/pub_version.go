package pubmodel

import (
	"private-pub-repo/modules/user/usermodel"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PubVersionModel struct {
	PackageName        string               `json:"package_name" gorm:"not null;index:,unique,composite:pubversion;"`
	Version            string               `json:"version" gorm:"not null;index:,unique,composite:pubversion;"`
	VersionNumberMajor uint64               `json:"version_number_major" gorm:"not null;"`
	VersionNumberMinor uint64               `json:"version_number_minor" gorm:"not null;"`
	VersionNumberPatch uint64               `json:"version_number_patch" gorm:"not null;"`
	Prerelease         bool                 `json:"prerelease" gorm:"not null;default:false;"`
	Pubspec            datatypes.JSON       `json:"pubspec" gorm:"not null;default:'{}';"`
	UploaderID         *uuid.UUID           `json:"user_id" gorm:"type:uuid;nullable;"`
	Uploader           *usermodel.UserModel `gorm:"foreignKey:UploaderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Readme             *string              `json:"readme" gorm:"type:text;nullable;"`
	Changelog          *string              `json:"changelog" gorm:"type:text;nullable;"`
	CreatedAt          *time.Time           `json:"created_at,omitempty" gorm:"not null;"`
	UpdatedAt          *time.Time           `json:"updated_at,omitempty" gorm:"not null;"`
	DeletedAt          *gorm.DeletedAt      `json:"deleted_at,omitempty" gorm:"index"`
}

func (PubVersionModel) TableName() string {
	return "pub_versions"
}
