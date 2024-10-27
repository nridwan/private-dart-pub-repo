package pubtokendto

import (
	"private-pub-repo/modules/pubtoken/pubtokenmodel"
	"time"

	"github.com/google/uuid"
)

type CreateTokenDTO struct {
	Remarks   string `json:"remarks" validate:"required,min=1"`
	Write     bool   `json:"write" validate:"required,boolean"`
	ExpiredAt string `json:"expired_at" validate:"required,datetime=2006-01-02"`
}

func (dto *CreateTokenDTO) ToModel(userId uuid.UUID, isAdmin bool) *pubtokenmodel.PubTokenModel {
	expiredAt, _ := time.Parse("2006-01-02", dto.ExpiredAt)
	expiredAt = time.Date(expiredAt.Year(), expiredAt.Month(), expiredAt.Day(), 23, 59, 59, 0, expiredAt.Location())

	return &pubtokenmodel.PubTokenModel{
		Remarks:   dto.Remarks,
		Write:     isAdmin && dto.Write,
		ExpiredAt: &expiredAt,
		UserID:    &userId,
	}
}
