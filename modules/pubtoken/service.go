package pubtoken

import (
	"context"
	"private-pub-repo/modules/app/appmodel"
	"private-pub-repo/modules/db"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/pubtoken/pubtokendto"
	"private-pub-repo/modules/pubtoken/pubtokenmodel"
	"private-pub-repo/utils"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	JwtIssuer = "pubToken"
)

type PubTokenService interface {
	Init(db db.DbService)
	Insert(context context.Context, pubtoken *pubtokenmodel.PubTokenModel) (*string, error)
	List(context context.Context, req *appmodel.GetListRequest, userId *uuid.UUID) (*appmodel.PaginationResponseList, error)
	Detail(context context.Context, id uuid.UUID, userId *uuid.UUID) (*pubtokenmodel.PubTokenModel, error)
	Update(context context.Context, id uuid.UUID, userId *uuid.UUID, updateDTO *pubtokendto.UpdateTokenDTO) (*pubtokenmodel.PubTokenModel, error)
	Delete(context context.Context, id uuid.UUID, userId *uuid.UUID) error
}

type pubTokenServiceImpl struct {
	monitorService monitor.MonitorService
	jwtService     jwt.JwtService
	db             *gorm.DB
}

func NewPubTokenService(jwtService jwt.JwtService, monitorService monitor.MonitorService) PubTokenService {
	return &pubTokenServiceImpl{
		jwtService:     jwtService,
		monitorService: monitorService,
	}
}

// impl `PubTokenService` start

func (service *pubTokenServiceImpl) Init(db db.DbService) {
	service.db = db.Default()
}

func (service *pubTokenServiceImpl) Insert(context context.Context, pubToken *pubtokenmodel.PubTokenModel) (*string, error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "PubTokenService.Insert", map[string]interface{}{})
	defer span.End()

	result := service.db.WithContext(spanContext).Create(pubToken)

	if result.Error != nil {
		return nil, result.Error
	}

	response, err := service.jwtService.GenerateAccessTokenTimed(pubToken.ID, JwtIssuer, time.Now().Unix(), map[string]interface{}{}, pubToken.ExpiredAt)
	return &response, err
}

func (service *pubTokenServiceImpl) List(context context.Context, req *appmodel.GetListRequest, userId *uuid.UUID) (*appmodel.PaginationResponseList, error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "PubTokenService.List", utils.StructToMap(req))
	defer span.End()
	var count int64
	pubtokens := []pubtokenmodel.PubTokenModel{}
	query := service.db.WithContext(spanContext).Model(pubtokens).Where("user_id = ?", userId)
	if req.Search != "" {
		query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Perform count and find concurrently using goroutines
	errChan := make(chan error, 2)
	go func() {
		defer wg.Done()
		errChan <- query.Session(&gorm.Session{}).Count(&count).Error
	}()

	go func() {
		defer wg.Done()
		query = query.Session(&gorm.Session{})
		errChan <- query.Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Find(&pubtokens).Error
	}()

	wg.Wait()

	var err error
	for i := 0; i < 2; i++ {
		select {
		case err = <-errChan:
			if err != nil {
				return nil, err
			}
		default:
		}
	}

	count32 := int(count)

	return &appmodel.PaginationResponseList{
		Pagination: &appmodel.PaginationResponsePagination{
			Page:  &req.Page,
			Size:  &req.Limit,
			Total: &count32,
		},
		Content: pubtokens,
	}, nil
}

func (service *pubTokenServiceImpl) Detail(context context.Context, id uuid.UUID, userId *uuid.UUID) (*pubtokenmodel.PubTokenModel, error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "PubTokenService.Detail", map[string]interface{}{
		"id": id.String(),
	})
	defer span.End()
	var pubToken pubtokenmodel.PubTokenModel
	result := service.db.WithContext(spanContext).First(&pubToken, pubtokendto.QueryTokenDTO{
		ID:     &id,
		UserID: userId,
	})

	return &pubToken, result.Error
}

func (service *pubTokenServiceImpl) Update(context context.Context, id uuid.UUID, userId *uuid.UUID, updateDTO *pubtokendto.UpdateTokenDTO) (*pubtokenmodel.PubTokenModel, error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "PubTokenService.Update", map[string]interface{}{
		"id": id.String(),
	})
	defer span.End()
	var pubToken pubtokenmodel.PubTokenModel
	result := service.db.WithContext(spanContext).Model(&pubToken).Where("id = ?", id).Where("user_id = ?", userId).Updates(updateDTO)

	if result.Error != nil {
		return nil, result.Error
	}

	return service.Detail(context, id, userId)
}

func (service *pubTokenServiceImpl) Delete(context context.Context, id uuid.UUID, userId *uuid.UUID) error {
	spanContext, span := service.monitorService.StartTraceSpan(context, "PubTokenService.Delete", map[string]interface{}{
		"id": id.String(),
	})
	defer span.End()
	var pubtoken pubtokenmodel.PubTokenModel
	result := service.db.WithContext(spanContext).Delete(&pubtoken, pubtokendto.QueryTokenDTO{
		ID:     &id,
		UserID: userId,
	})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

// impl `PubTokenService` end
