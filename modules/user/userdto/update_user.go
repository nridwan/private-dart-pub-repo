package userdto

type UpdateUserDTO struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=4"`
	IsAdmin  *bool   `json:"is_admin" validate:"omitempty"`
	CanWrite *bool   `json:"can_write" validate:"omitempty"`
}
