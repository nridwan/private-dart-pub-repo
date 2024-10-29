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
	"private-pub-repo/modules/storage"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const filePathFormat = "pub/packages/%s/versions/%s.tar.gz"

type PubService interface {
	Init(db db.DbService)
	VersionList(context context.Context, packageName string, baseUrl string, publicOnly bool) (*pubdto.PubPackageDTO, error)
	VersionDetail(context context.Context, packageName string, version string, baseUrl string, publicOnly bool) (*pubdto.PubVersionDTO, error)
	GetUpstreamUrl(context context.Context, path string) *string
	UploadVersion(context context.Context, file *multipart.FileHeader, userId uuid.UUID) error
	GetDownloadUrl(context context.Context, packageName string, version string, baseUrl string, publicOnly bool) (*string, error)
}

type pubServiceImpl struct {
	monitorService monitor.MonitorService
	jwtService     jwt.JwtService
	db             *gorm.DB
	upstreamUrl    string
	storage        storage.StorageService
}

func NewPubService(jwtService jwt.JwtService, monitorService monitor.MonitorService, config *config.ConfigModule, storage storage.StorageService) PubService {
	return &pubServiceImpl{
		jwtService:     jwtService,
		monitorService: monitorService,
		upstreamUrl:    config.Getenv("UPSTREAM_URL", ""),
		storage:        storage,
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

func (service *pubServiceImpl) GetUpstreamUrl(context context.Context, path string) *string {
	_, span := service.monitorService.StartTraceSpan(context, "PubService.GetUpstreamUrl", map[string]interface{}{})
	defer span.End()
	if service.upstreamUrl == "" {
		return nil
	}
	newUrl := service.upstreamUrl + path
	return &newUrl
}

func (service *pubServiceImpl) UploadVersion(context context.Context, file *multipart.FileHeader, userId uuid.UUID) error {
	_, span := service.monitorService.StartTraceSpan(context, "PubService.UploadVersion", map[string]interface{}{})
	defer span.End()
	tarPackageInfo := pubdto.TarPackageInfoDTO{}

	// Loop through each entry in the tar archive
	hasPubspec, shouldReturn, returnValue := service.readArchiveContent(file, &tarPackageInfo)
	if shouldReturn {
		return returnValue
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

	err = service.storage.Upload(fmt.Sprintf(filePathFormat, packageName, version), file)

	if err != nil {
		return err
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

func (service *pubServiceImpl) readArchiveContent(file *multipart.FileHeader, tarPackageInfo *pubdto.TarPackageInfoDTO) (bool, bool, error) {
	reader, err := file.Open()
	if err != nil {
		return false, true, err
	}
	defer reader.Close()

	// Create a Gzip reader
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return false, true, err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	hasPubspec := false

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, true, err
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
					return false, true, err
				}

				var data map[string]interface{}
				if err := yaml.Unmarshal(content, &data); err != nil {
					return false, true, fmt.Errorf("failed to unmarshal pubspec.yaml: %w", err)
				}
				tarPackageInfo.Pubspec = data
				hasPubspec = true
			case "readme.md":
				content, err = io.ReadAll(tarReader)
				if err != nil {
					return false, true, err
				}
				tarPackageInfo.Readme = string(content)
			case "changelog.md":
				content, err = io.ReadAll(tarReader)
				if err != nil {
					return false, true, err
				}
				tarPackageInfo.Changelog = string(content)
			}
		}
	}
	return hasPubspec, false, nil
}

func (service *pubServiceImpl) GetDownloadUrl(context context.Context, packageName string, version string, baseUrl string, publicOnly bool) (*string, error) {
	_, span := service.monitorService.StartTraceSpan(context, "PubService.GetDownloadUrl", map[string]interface{}{})
	defer span.End()

	_, err := service.VersionDetail(context, packageName, version, baseUrl, publicOnly)

	if err != nil {
		return nil, err
	}

	url := service.storage.GetUrl(fmt.Sprintf(filePathFormat, packageName, version))
	return &url, nil
}

// impl `PubService` end
