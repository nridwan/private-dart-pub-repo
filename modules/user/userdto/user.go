package userdto

import (
	"private-pub-repo/modules/user/usermodel"
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID        uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;primaryKey;default:uuid_generate_v4()"`
	Name      string     `json:"name" gorm:"not null;"`
	Email     string     `json:"email" gorm:"not null;unique;"`
	IsAdmin   bool       `json:"is_admin" gorm:"not null;unique;"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"not null;"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"not null;"`
}

func MapUserModelToDTO(model *usermodel.UserModel) *UserDTO {
	return &UserDTO{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		IsAdmin:   model.IsAdmin,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
