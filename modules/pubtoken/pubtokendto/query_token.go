package pubtokendto

import "github.com/google/uuid"

type QueryTokenDTO struct {
	ID     *uuid.UUID `json:"id"`
	UserID *uuid.UUID `json:"user_id"`
}
