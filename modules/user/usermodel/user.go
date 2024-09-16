package usermodel

import "private-pub-repo/base"

type UserModel struct {
	base.BaseModel
	Name     string  `json:"name" gorm:"not null;"`
	Email    string  `json:"email" gorm:"not null;unique;"`
	Password *string `json:"-" gorm:"not null;"`
	IsAdmin  bool    `json:"is_admin" gorm:"not null;default:false"`
}

func (UserModel) TableName() string {
	return "users"
}
