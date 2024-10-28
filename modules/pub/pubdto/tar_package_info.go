package pubdto

type TarPackageInfoDTO struct {
	Changelog string
	Readme    string
	Pubspec   map[string]interface{}
}
