package pubtokendto

import (
	"private-pub-repo/modules/pubtoken/pubtokenmodel"
	"time"

	"github.com/google/uuid"
)

type CreateTokenDTO struct {
	Remarks   string `json:"remarks" validate:"required,min=1"`
	Read      bool   `json:"read" validate:"required,boolean"`
	Write     bool   `json:"write" validate:"required,boolean"`
	ExpiredAt string `json:"expired_at" validate:"required,datetime=2006-01-02"`
}

func (dto *CreateTokenDTO) ToModel(userId uuid.UUID) *pubtokenmodel.PubTokenModel {
	expiredAt, _ := time.Parse("2006-01-02", dto.ExpiredAt)

	return &pubtokenmodel.PubTokenModel{
		Remarks:   dto.Remarks,
		Read:      dto.Read,
		Write:     dto.Write,
		ExpiredAt: &expiredAt,
		UserID:    userId,
	}
}
