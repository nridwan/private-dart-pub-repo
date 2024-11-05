package user

import (
	"context"
	"fmt"
	"private-pub-repo/base"
	"private-pub-repo/modules/app/appmodel"
	"private-pub-repo/modules/config"
	"private-pub-repo/modules/db"
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/mail"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/user/userdto"
	"private-pub-repo/modules/user/usermodel"
	"private-pub-repo/utils"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	jwtIssuer = "appUser"

	emailWhereQuery = "email = ?"
)

type UserService interface {
	Init(db db.DbService)
	Insert(context context.Context, user *usermodel.UserModel) (*userdto.UserDTO, error)
	Update(context context.Context, id uuid.UUID, updateDTO *userdto.UpdateUserDTO) (*userdto.UserDTO, error)
	List(context context.Context, req *appmodel.GetListRequest) (*appmodel.PaginationResponseList, error)
	Detail(context context.Context, id uuid.UUID) (*userdto.UserDTO, error)
	Delete(context context.Context, id uuid.UUID) error
	Login(context context.Context, req *userdto.LoginDTO) (*userdto.LoginResponseDTO, error)
	ForgotOtp(context context.Context, req *userdto.ForgotOtpDTO) (*userdto.ForgotOtpResponseDTO, error)
	ForgotCreatePassword(context context.Context, req *userdto.ForgotCreatePasswordDTO) (bool, error)
	RefreshToken(context context.Context, claims jwt.JwtClaim) (response *userdto.LoginResponseDTO, err error)
	GenerateHashPassword(password string) (*string, error)
}

type userServiceImpl struct {
	monitorService monitor.MonitorService
	jwtService     jwt.JwtService
	db             *gorm.DB
	otpExpiredTime time.Duration
	mail           mail.MailService
}

func NewUserService(jwtService jwt.JwtService, monitorService monitor.MonitorService, config config.ConfigService, mail mail.MailService) UserService {
	otpExpiredTime, err := strconv.Atoi(config.Getenv("OTP_EXPIRED_TIME", "5"))

	if err != nil {
		otpExpiredTime = 5
	}

	return &userServiceImpl{
		jwtService:     jwtService,
		monitorService: monitorService,
		otpExpiredTime: time.Duration(otpExpiredTime) * time.Minute,
		mail:           mail,
	}
}

func (service *userServiceImpl) validateEmail(context context.Context, email string) error {
	var count int64
	service.db.WithContext(context).Model(&usermodel.UserModel{}).Where(emailWhereQuery, email).Count(&count)
	if count > 0 {
		return fiber.NewError(400, "Email already registered")
	}
	return nil
}

// impl `UserService` start

func (service *userServiceImpl) Init(db db.DbService) {
	service.db = db.Default()
}

func (service *userServiceImpl) Insert(context context.Context, user *usermodel.UserModel) (*userdto.UserDTO, error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "UserService.Insert", map[string]interface{}{})
	defer span.End()
	err := service.validateEmail(spanContext, user.Email)
	if err != nil {
		return nil, err
	}

	pwd, err := service.GenerateHashPassword(*user.Password)

	if err != nil {
		return nil, err
	}

	user.Password = pwd
	result := service.db.WithContext(spanContext).Create(user)
	dto := userdto.MapUserModelToDTO(user)
	dto.UpdatedAt = nil
	return dto, result.Error
}

func (service *userServiceImpl) Update(context context.Context, id uuid.UUID, updateDTO *userdto.UpdateUserDTO) (*userdto.UserDTO, error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "UserService.Update", map[string]interface{}{
		"id": id.String(),
	})
	defer span.End()
	if updateDTO.Password != nil {
		updateDTO.Password = nil
	}
	user := usermodel.UserModel{BaseModel: base.BaseModel{ID: id}}
	result := service.db.WithContext(spanContext).Model(&user).Updates(updateDTO)
	if result.Error != nil {
		return nil, result.Error
	}
	return service.Detail(context, id)
}

func (service *userServiceImpl) List(context context.Context, req *appmodel.GetListRequest) (*appmodel.PaginationResponseList, error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "UserService.List", utils.StructToMap(req))
	defer span.End()
	var count int64
	users := []usermodel.UserModel{}
	query := service.db.WithContext(spanContext).Model(users)
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
		errChan <- query.Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Find(&users).Error
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
		Content: users,
	}, nil
}

func (service *userServiceImpl) Detail(context context.Context, id uuid.UUID) (*userdto.UserDTO, error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "UserService.Detail", map[string]interface{}{
		"id": id.String(),
	})
	defer span.End()
	var user usermodel.UserModel
	result := service.db.WithContext(spanContext).First(&user, id)
	return userdto.MapUserModelToDTO(&user), result.Error
}

func (service *userServiceImpl) Delete(context context.Context, id uuid.UUID) error {
	spanContext, span := service.monitorService.StartTraceSpan(context, "UserService.Delete", map[string]interface{}{
		"id": id.String(),
	})
	defer span.End()
	var user userdto.UserDTO
	result := service.db.WithContext(spanContext).Delete(&user, id)
	return result.Error
}

func (service *userServiceImpl) Login(context context.Context, req *userdto.LoginDTO) (response *userdto.LoginResponseDTO, err error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "UserService.Login", map[string]interface{}{
		"email": req.Email,
	})
	defer span.End()
	var user usermodel.UserModel
	result := service.db.WithContext(spanContext).Where(emailWhereQuery, req.Email).First(&user)
	if result.Error != nil {
		err = fiber.NewError(400, "Email and password doesn't match.")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)) != nil {
		return nil, fiber.NewError(400, "Email and password doesn't match.")
	}

	response, err = service.jwtService.GenerateToken(user.ID, jwtIssuer, map[string]interface {
	}{
		"is_admin": user.IsAdmin,
	})
	return
}

func (service *userServiceImpl) ForgotOtp(context context.Context, req *userdto.ForgotOtpDTO) (response *userdto.ForgotOtpResponseDTO, err error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "UserService.ForgotOtp", map[string]interface{}{
		"email": req.Email,
	})
	defer span.End()
	var user usermodel.UserModel
	result := service.db.WithContext(spanContext).Where(emailWhereQuery, req.Email).First(&user)
	if result.Error != nil {
		err = result.Error
		return
	}

	otp, err := gonanoid.Generate("1234567890", 6)

	if err != nil {
		return
	}

	hashOtp, err := service.GenerateHashPassword(otp)

	if err != nil {
		return
	}

	expiredAt := time.Now().Add(service.otpExpiredTime)

	userOtp := usermodel.UserOtpModel{
		ID:        user.ID,
		Purpose:   usermodel.OtpPurposeForgot,
		Otp:       *hashOtp,
		ExpiredAt: &expiredAt,
	}

	result = service.db.WithContext(spanContext).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"purpose", "otp", "expired_at"}),
	}).Create(userOtp)
	if result.Error != nil {
		err = result.Error
		return
	}

	err = service.mail.Send(
		[]string{user.Email},
		[]string{},
		"Forgot Password OTP",
		fmt.Sprintf(
			`Please use this OTP below to change your password: %s. 
OTP only valid for %d minute.
Do not share this to anyone.`,
			otp, service.otpExpiredTime/time.Minute,
		),
	)
	if err != nil {
		return
	}

	response = &otp

	return
}

func (service *userServiceImpl) ForgotCreatePassword(context context.Context, req *userdto.ForgotCreatePasswordDTO) (response bool, err error) {
	spanContext, span := service.monitorService.StartTraceSpan(context, "UserService.ForgotCreatePassword", map[string]interface{}{
		"email": req.Email,
	})
	defer span.End()
	var user usermodel.UserModel
	result := service.db.WithContext(spanContext).Model(user).
		Select("\"users\".\"id\"", "\"UserOtp\".\"otp\"").
		InnerJoins("UserOtp", service.db.Where("purpose = ?", usermodel.OtpPurposeForgot)).
		Where(emailWhereQuery, req.Email).
		Where("expired_at >= NOW()").
		First(&user)
	if result.Error != nil {
		err = result.Error
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.UserOtp.Otp), []byte(req.Otp)) != nil {
		err = fiber.ErrNotFound
		return
	}

	pwd, err := service.GenerateHashPassword(req.Password)

	if err != nil {
		return
	}

	result = service.db.WithContext(spanContext).Model(&user).Updates(userdto.UpdateUserDTO{
		Password: pwd,
	})

	if result.Error != nil {
		err = result.Error
		return
	}

	service.db.WithContext(spanContext).Delete(&usermodel.UserOtpModel{}, user.ID)

	response = true
	return
}

func (service *userServiceImpl) RefreshToken(context context.Context, claims jwt.JwtClaim) (response *userdto.LoginResponseDTO, err error) {
	_, span := service.monitorService.StartTraceSpan(context, "UserService.RefreshToken", map[string]interface{}{})
	defer span.End()
	response, err = service.jwtService.Refresh(claims)
	return
}

func (service *userServiceImpl) GenerateHashPassword(password string) (*string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	pwdString := string(pwd)

	return &pwdString, err
}

// impl `UserService` end
