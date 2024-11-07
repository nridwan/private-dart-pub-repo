package userdto

import (
	"private-pub-repo/modules/user/usermodel"
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID        uuid.UUID  `json:"user_id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	IsAdmin   bool       `json:"is_admin"`
	CanWrite  bool       `json:"can_write"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func MapUserModelToDTO(model *usermodel.UserModel) *UserDTO {
	return &UserDTO{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		IsAdmin:   model.IsAdmin,
		CanWrite:  model.CanWrite,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
