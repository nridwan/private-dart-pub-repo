package userdto

type ForgotCreatePasswordDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Otp      string `json:"otp" validate:"required,len=6"`
	Password string `json:"password" validate:"required,min=4"`
}
