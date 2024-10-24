package user

import (
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/utils"

	"github.com/gofiber/fiber/v2"
)

type UserJwtMiddleware interface {
	jwt.JwtMiddleware
	IsAdmin(c *fiber.Ctx) error
}

type userMiddlewareImpl struct {
	jwtService     jwt.JwtService
	monitorService monitor.MonitorService
}

func NewUserJwtMiddleware(jwtService jwt.JwtService, monitorService monitor.MonitorService) UserJwtMiddleware {
	return &userMiddlewareImpl{
		jwtService:     jwtService,
		monitorService: monitorService,
	}
}

// impl `UserJwtMiddleware` start

func (service *userMiddlewareImpl) IsAdmin(c *fiber.Ctx) error {
	if utils.GetFiberJwtClaims(c)["is_admin"] != true {
		return fiber.NewError(401, "Unauthenticated")
	}

	return c.Next()
}

// impl `UserJwtMiddleware` end

// impl `jwt.JwtMiddleware` start

func (service *userMiddlewareImpl) CanAccess(c *fiber.Ctx) error {
	err := service.jwtService.CanAccess(c, jwtIssuer)

	if err == nil {
		var userId string
		if userId, err = utils.GetFiberJwtUserIdString(c); err == nil {
			service.monitorService.SetCurrentSpanAttributes(c.UserContext(), map[string]interface{}{"admin_user_id": userId})
		}
	}

	return err
}
func (service *userMiddlewareImpl) CanRefresh(c *fiber.Ctx) error {
	return service.jwtService.CanRefresh(c, jwtIssuer)
}

// impl `jwt.JwtMiddleware` end
