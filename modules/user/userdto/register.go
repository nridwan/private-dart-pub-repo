package userdto

import "private-pub-repo/modules/user/usermodel"

type RegisterDTO struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=4"`
	Email    string `json:"email" validate:"required,email"`
	CanWrite bool   `json:"can_write" validate:"boolean"`
	IsAdmin  bool   `json:"is_admin" validate:"boolean"`
}

func (dto *RegisterDTO) ToModel() *usermodel.UserModel {
	return &usermodel.UserModel{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: &dto.Password,
		CanWrite: dto.CanWrite,
		IsAdmin:  dto.IsAdmin,
	}
}
