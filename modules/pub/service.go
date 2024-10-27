package pub

import (
	"context"
	"private-pub-repo/modules/config"
	"private-pub-repo/modules/db"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/pub/pubdto"
	"private-pub-repo/modules/pub/pubmodel"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PubService interface {
	Init(db db.DbService)
	VersionList(context context.Context, packageName string, baseUrl string, publicOnly bool) (*pubdto.PubPackageDTO, error)
	VersionDetail(context context.Context, packageName string, version string, baseUrl string, publicOnly bool) (*pubdto.PubVersionDTO, error)
	GetUpstreamUrl(path string) *string
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

// impl `PubService` end
