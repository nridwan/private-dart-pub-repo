package pub

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"private-pub-repo/modules/config"
	"private-pub-repo/modules/db"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/pub/pubdto"
	"private-pub-repo/modules/pub/pubmodel"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PubService interface {
	Init(db db.DbService)
	VersionList(context context.Context, packageName string, baseUrl string, publicOnly bool) (*pubdto.PubPackageDTO, error)
	VersionDetail(context context.Context, packageName string, version string, baseUrl string, publicOnly bool) (*pubdto.PubVersionDTO, error)
	GetUpstreamUrl(path string) *string
	UploadVersion(file *multipart.FileHeader, userId uuid.UUID) error
}

type pubServiceImpl struct {
	monitorService monitor.MonitorService
	jwtService     jwt.JwtService
	db             *gorm.DB
	upstreamUrl    string
}

func NewPubService(jwtService jwt.JwtService, monitorService monitor.MonitorService, config *config.ConfigModule) PubService {
	return &pubServiceImpl{
		jwtService:     jwtService,
		monitorService: monitorService,
		upstreamUrl:    config.Getenv("UPSTREAM_URL", ""),
	}
}

// impl `PubService` start

func (service *pubServiceImpl) Init(db db.DbService) {
	service.db = db.Default()
}

func (service *pubServiceImpl) VersionList(context context.Context, packageName string, baseUrl string, publicOnly bool) (*pubdto.PubPackageDTO, error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "PubService.VersionList", map[string]interface{}{})
	defer span.End()
	pubPackage := pubmodel.PubPackageModel{}
	pubVersions := []pubmodel.PubVersionModel{}

	result := service.db.WithContext(spanContext).First(&pubPackage, "name = ?", packageName)

	if result.Error != nil {
		return nil, fiber.ErrNotFound
	}

	if pubPackage.Private && publicOnly {
		return nil, fiber.ErrForbidden
	}

	service.db.WithContext(spanContext).Model(pubVersions).
		Select("package_name", "version", "pubspec").
		Where("package_name = ?", packageName).
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "version_number_major"}, Desc: true},
			{Column: clause.Column{Name: "version_number_minor"}, Desc: true},
			{Column: clause.Column{Name: "version_number_patch"}, Desc: true},
		}}).Find(&pubVersions)

	if len(pubVersions) > 0 {
		pubDTO := pubdto.MapPubVersionsToPackageDTO(pubVersions, baseUrl)
		return &pubDTO, nil
	}

	return nil, fiber.ErrNotFound
}

func (service *pubServiceImpl) VersionDetail(context context.Context, packageName string, version string, baseUrl string, publicOnly bool) (*pubdto.PubVersionDTO, error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "PubService.VersionDetail", map[string]interface{}{})
	defer span.End()
	pubPackage := pubmodel.PubPackageModel{}
	pubVersion := pubmodel.PubVersionModel{}

	result := service.db.WithContext(spanContext).First(&pubPackage, "name = ?", packageName)

	if result.Error != nil {
		return nil, fiber.ErrNotFound
	}

	if pubPackage.Private && publicOnly {
		return nil, fiber.ErrForbidden
	}

	result = service.db.WithContext(spanContext).Model(pubVersion).
		Select("package_name", "version", "pubspec").
		Where("package_name = ?", packageName).
		Where("version = ?", version).
		First(&pubVersion)

	if result.Error != nil {
		return nil, fiber.ErrNotFound
	}

	pubDTO := pubdto.MapPubVersionToDTO(&pubVersion, baseUrl)

	return &pubDTO, nil
}

func (service *pubServiceImpl) GetUpstreamUrl(path string) *string {
	if service.upstreamUrl == "" {
		return nil
	}
	newUrl := service.upstreamUrl + path
	return &newUrl
}

func (service *pubServiceImpl) UploadVersion(file *multipart.FileHeader, userId uuid.UUID) error {
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create a Gzip reader
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	tarPackageInfo := pubdto.TarPackageInfoDTO{}

	hasPubspec := false

	// Loop through each entry in the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Check if the filename matches any target file (case insensitive)
		filename := filepath.Base(header.Name)

		// Extract the file content based on its type
		switch header.Typeflag {
		case tar.TypeReg:
			var content []byte

			// Extract pubspec.yaml content into map
			switch strings.ToLower(filename) {
			case "pubspec.yaml":
				content, err = io.ReadAll(tarReader)
				if err != nil {
					return err
				}

				var data map[string]interface{}
				if err := yaml.Unmarshal(content, &data); err != nil {
					return fmt.Errorf("failed to unmarshal pubspec.yaml: %w", err)
				}
				tarPackageInfo.Pubspec = data
				hasPubspec = true
			case "readme.md":
				content, err = io.ReadAll(tarReader)
				if err != nil {
					return err
				}
				tarPackageInfo.Readme = string(content)
			case "changelog.md":
				content, err = io.ReadAll(tarReader)
				if err != nil {
					return err
				}
				tarPackageInfo.Changelog = string(content)
			}
		}
	}

	if !hasPubspec {
		return fmt.Errorf("did not find any pubspec.yaml file in upload, aborting")
	}

	parseOk := true

	packageName, ok := tarPackageInfo.Pubspec["name"].(string)
	parseOk = parseOk && ok

	version, ok := tarPackageInfo.Pubspec["version"].(string)
	parseOk = parseOk && ok

	pubspecJson, err := json.Marshal(tarPackageInfo.Pubspec)

	semverObj, errSemver := semver.NewVersion(version)

	if !parseOk || err != nil || errSemver != nil {
		return fmt.Errorf("invalid pubspec.yaml")
	}

	result := service.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).Create(&pubmodel.PubPackageModel{Name: packageName})

	if result.Error != nil {
		return result.Error
	}

	result = service.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "package_name"}, {Name: "version"}},
		DoUpdates: clause.AssignmentColumns([]string{"readme", "changelog", "pubspec", "uploader_id"}),
	}).Create(&pubmodel.PubVersionModel{
		PackageName: packageName, Version: version,
		VersionNumberMajor: semverObj.Major(),
		VersionNumberMinor: semverObj.Minor(),
		VersionNumberPatch: semverObj.Patch(),
		Prerelease:         semverObj.Prerelease() != "",
		Readme:             &tarPackageInfo.Readme,
		Changelog:          &tarPackageInfo.Changelog,
		Pubspec:            pubspecJson,
		UploaderID:         &userId,
	})

	if result.Error != nil {
		return result.Error
	}

	//upload tar.gz
	return nil
}

// impl `PubService` end
