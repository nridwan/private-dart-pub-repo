package usermodel

import (
	"github.com/google/uuid"
)

type UserOtpPurpose = string

const OtpPurposeForgot UserOtpPurpose = "forgot"

type UserOtpModel struct {
	ID      uuid.UUID      `json:"id" gorm:"type:uuid;not null;primaryKey"`
	Purpose UserOtpPurpose `json:"purpose" gorm:"not null;"`
	Otp     string         `json:"-" gorm:"not null;"`
}

func (UserOtpModel) TableName() string {
	return "user_otps"
}
