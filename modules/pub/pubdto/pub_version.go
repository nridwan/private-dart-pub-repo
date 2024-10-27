package pubdto

import (
	"encoding/json"
	"fmt"
	"private-pub-repo/modules/pub/pubmodel"
)

type PubVersionDTO struct {
	Version    string                 `json:"version"`
	ArchiveUrl string                 `json:"archive_url"`
	Pubspec    map[string]interface{} `json:"pubspec"`
}

func MapPubVersionToDTO(model *pubmodel.PubVersionModel, baseUrl string) PubVersionDTO {
	archiveUrl := fmt.Sprintf("%s/v1/pub/packages/%s/versions/%s.tar.gz", baseUrl, model.PackageName, model.Version)

	var pubspec map[string]interface{}
	json.Unmarshal([]byte(model.Pubspec), &pubspec) // Convert JSON to map

	return PubVersionDTO{
		Version:    model.Version,
		ArchiveUrl: archiveUrl,
		Pubspec:    pubspec,
	}
}
