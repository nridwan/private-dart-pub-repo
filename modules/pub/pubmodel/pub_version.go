package pubmodel

import (
	"private-pub-repo/modules/user/usermodel"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PubVersionModel struct {
	PackageName        string                 `json:"package_name" gorm:"not null;"`
	Version            string                 `json:"version" gorm:"not null;primaryKey;"`
	VersionNumberMajor uint                   `json:"version_number_major" gorm:"not null;"`
	VersionNumberMinor uint                   `json:"version_number_minor" gorm:"not null;"`
	VersionNumberPatch uint                   `json:"version_number_patch" gorm:"not null;"`
	Prerelease         bool                   `json:"prerelease" gorm:"not null;default:false;"`
	Remarks            string                 `json:"remarks" gorm:"not null;"`
	Pubspec            map[string]interface{} `json:"pubspec" gorm:"not null;default:'{}';"`
	UploaderID         *uuid.UUID             `json:"user_id" gorm:"type:uuid;nullable;"`
	Uploader           *usermodel.UserModel   `gorm:"foreignKey:UploaderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Readme             *string                `json:"readme" gorm:"type:text;nullable;"`
	Changelog          *string                `json:"changelog" gorm:"type:text;nullable;"`
	CreatedAt          *time.Time             `json:"created_at,omitempty" gorm:"not null;"`
	UpdatedAt          *time.Time             `json:"updated_at,omitempty" gorm:"not null;"`
	DeletedAt          *gorm.DeletedAt        `json:"deleted_at,omitempty" gorm:"index"`
}

func (PubVersionModel) TableName() string {
	return "pub_versions"
}
