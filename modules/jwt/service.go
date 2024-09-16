package jwt

import (
	"private-pub-repo/modules/config"
	"strconv"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	JwtAppAud     = "app"
	JwtRefreshAud = "token_refresh"
)

type JwtService interface {
	jwtCommonMiddleware
	Init(config config.ConfigService)
	GetSecret() string
	GetHandler() fiber.Handler
	Refresh(claims JwtClaim) (*JWTTokenModel, error)
	GenerateToken(id uuid.UUID, issuer string, payload map[string]interface{}) (*JWTTokenModel, error)
}

func (service *JwtModule) errorHandler(ctx *fiber.Ctx, err error) error {
	println(err.Error())
	return fiber.NewError(401, "Unauthenticated")
}

// impl `JwtService` start

func (service *JwtModule) Init(config config.ConfigService) {
	service.secret = config.Getenv("JWT_SECRET", "")
	service.handler = jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(service.secret)},
		ErrorHandler: service.errorHandler,
	})
	if localLifetime, err := strconv.Atoi(config.Getenv("JWT_TOKEN_LIFETIME", "1")); err == nil {
		service.lifetime = time.Duration(localLifetime)
	}
	if localLifetime, err := strconv.Atoi(config.Getenv("JWT_REFRESH_LIFETIME", "1")); err == nil {
		service.refreshLifetime = time.Duration(localLifetime)
	}
}

func (service *JwtModule) GetSecret() string {
	return service.secret
}

func (service *JwtModule) GetHandler() fiber.Handler {
	return service.handler
}

func (service *JwtModule) Refresh(claims JwtClaim) (*JWTTokenModel, error) {
	now := time.Now().Unix()

	sub, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		return nil, err
	}

	issuer := claims["iss"].(string)

	// Generate encoded token and send it as response.
	accessToken, err := service.generateAccessToken(sub, issuer, now, claims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := service.generateRefreshToken(sub, issuer, now, claims)
	if err != nil {
		return nil, err
	}

	return &JWTTokenModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (service *JwtModule) GenerateToken(id uuid.UUID, issuer string, payload map[string]interface{}) (*JWTTokenModel, error) {
	now := time.Now().Unix()
	// Generate encoded token and send it as response.
	accessToken, err := service.generateAccessToken(id, issuer, now, payload)
	if err != nil {
		return nil, err
	}
	refreshToken, err := service.generateRefreshToken(id, issuer, now, payload)
	if err != nil {
		return nil, err
	}

	return &JWTTokenModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (service *JwtModule) generateAccessToken(id uuid.UUID, issuer string, now int64, payload map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(JwtClaim)

	mergeJwtClaims(payload, claims)

	claims["sub"] = id
	claims["iat"] = now
	claims["nbf"] = now
	claims["exp"] = time.Unix(now, 0).Add(time.Minute * service.lifetime).Unix()
	claims["iss"] = issuer
	claims["aud"] = []string{JwtAppAud}
	// Generate encoded token and send it as response.
	return token.SignedString([]byte(service.GetSecret()))
}

func (service *JwtModule) generateRefreshToken(id uuid.UUID, issuer string, now int64, payload map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	currentTime := time.Now()
	// Set claims
	claims := token.Claims.(JwtClaim)

	mergeJwtClaims(payload, claims)

	claims["sub"] = id
	claims["iat"] = currentTime.Unix()
	claims["nbf"] = currentTime.Unix()
	claims["exp"] = time.Unix(now, 0).Add(time.Minute * service.refreshLifetime).Unix()
	claims["iss"] = issuer
	claims["aud"] = []string{JwtRefreshAud}
	// Generate encoded token and send it as response.
	return token.SignedString([]byte(service.GetSecret()))
}

// impl `JwtService` end
