package pubtokenmodel

import (
	"private-pub-repo/base"
	"private-pub-repo/modules/user/usermodel"
	"time"

	"github.com/google/uuid"
)

type PubTokenModel struct {
	base.BaseModel
	Remarks   string              `json:"remarks" gorm:"not null;"`
	Read      bool                `json:"read" gorm:"not null;default:true"`
	Write     bool                `json:"write" gorm:"not null;default:false"`
	ExpiredAt *time.Time          `json:"expired_at" gorm:"not null;"`
	UserID    uuid.UUID           `json:"user_id" gorm:"type:uuid;nullable;"`
	User      usermodel.UserModel `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (PubTokenModel) TableName() string {
	return "pub_tokens"
}
