package pubtokendto

type UpdateTokenDTO struct {
	Write *bool `json:"write" validate:"boolean"`
}
