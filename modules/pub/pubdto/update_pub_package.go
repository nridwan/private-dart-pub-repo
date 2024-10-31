package pubdto

type UpdatePubPackageDTO struct {
	Private bool `json:"private" validate:"required,boolean"`
}
