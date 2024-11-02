package usermodel

import (
	"time"

	"github.com/google/uuid"
)

type UserOtpPurpose = string

const OtpPurposeForgot UserOtpPurpose = "forgot"

type UserOtpModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;not null;primaryKey"`
	Purpose   UserOtpPurpose `json:"purpose" gorm:"not null;"`
	Otp       string         `json:"-" gorm:"not null;"`
	ExpiredAt *time.Time     `json:"expired_at,omitempty" gorm:"nullable;"`
}

func (UserOtpModel) TableName() string {
	return "user_otps"
}
