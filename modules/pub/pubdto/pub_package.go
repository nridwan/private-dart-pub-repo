package pubdto

import "private-pub-repo/modules/pub/pubmodel"

type PubPackageDTO struct {
	Name     string          `json:"name"`
	Latest   *PubVersionDTO  `json:"latest"`
	Versions []PubVersionDTO `json:"versions"`
}

func MapPubVersionsToPackageDTO(versions []pubmodel.PubVersionModel, baseUrl string) PubPackageDTO {
	var name string
	var latest *PubVersionDTO
	versionDTOs := make([]PubVersionDTO, len(versions))

	for i, version := range versions {
		versionDTOs[i] = MapPubVersionToDTO(&version, baseUrl)
		if !version.Prerelease {
			latest = &versionDTOs[i]
		}
	}

	if len(versions) > 0 {
		name = versions[0].PackageName
	}

	return PubPackageDTO{
		Name:     name,
		Latest:   latest,
		Versions: versionDTOs,
	}
}
