package userdto

type ForgotOtpDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotOtpResponseDTO = string
