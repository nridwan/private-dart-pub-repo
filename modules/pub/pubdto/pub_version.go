package pubdto

import (
	"fmt"
	"private-pub-repo/modules/pub/pubmodel"
)

type PubVersionDTO struct {
	Version    string                 `json:"version"`
	ArchiveUrl string                 `json:"archive_url"`
	Pubspec    map[string]interface{} `json:"pubspec"`
}

func MapPubVersionToDTO(model *pubmodel.PubVersionModel, baseUrl string) PubVersionDTO {
	archiveUrl := fmt.Sprintf("%s/packages/%s/versions/%s.tar.gz", baseUrl, model.PackageName, model.Version)
	return PubVersionDTO{
		Version:    model.Version,
		ArchiveUrl: archiveUrl,
		Pubspec:    model.Pubspec,
	}
}
