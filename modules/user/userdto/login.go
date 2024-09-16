package userdto

import (
	"gofiber-boilerplate/modules/jwt"
)

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponseDTO = jwt.JWTTokenModel
