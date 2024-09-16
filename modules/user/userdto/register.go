package userdto

import "gofiber-boilerplate/modules/user/usermodel"

type RegisterDTO struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=4"`
	Email    string `json:"email" validate:"required,email"`
}

func (dto *RegisterDTO) ToModel() *usermodel.UserModel {
	return &usermodel.UserModel{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: &dto.Password,
	}
}
