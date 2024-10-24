package pubtoken

import (
	"private-pub-repo/modules/jwt"
	"private-pub-repo/modules/monitor"
	"private-pub-repo/modules/pubtoken/pubtokenmodel"
	"private-pub-repo/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PubTokenJwtMiddleware interface {
	jwt.JwtMiddleware
	CanRead(c *fiber.Ctx) error
	CanWrite(c *fiber.Ctx) error
}

type pubTokenMiddlewareImpl struct {
	jwtService      jwt.JwtService
	pubTokenService PubTokenService
	monitorService  monitor.MonitorService
}

func NewUserJwtMiddleware(jwtService jwt.JwtService, pubTokenService PubTokenService, monitorService monitor.MonitorService) PubTokenJwtMiddleware {
	return &pubTokenMiddlewareImpl{
		jwtService:      jwtService,
		pubTokenService: pubTokenService,
		monitorService:  monitorService,
	}
}

// impl `UserJwtMiddleware` start

func (service *pubTokenMiddlewareImpl) CanRead(c *fiber.Ctx) error {
	if c.Locals("read") != true {
		return fiber.NewError(401, "Unauthenticated")
	}

	return c.Next()
}

func (service *pubTokenMiddlewareImpl) CanWrite(c *fiber.Ctx) error {
	if c.Locals("write") != true {
		return fiber.NewError(401, "Unauthenticated")
	}

	return c.Next()
}

// impl `UserJwtMiddleware` end

// impl `jwt.JwtMiddleware` start

func (service *pubTokenMiddlewareImpl) CanAccess(c *fiber.Ctx) error {
	err := service.jwtService.CanAccess(c, jwtIssuer)

	if err == nil {
		var pubTokenIdString string
		pubTokenIdString, err = utils.GetFiberJwtUserIdString(c)
		if err == nil {
			service.monitorService.SetCurrentSpanAttributes(c.UserContext(), map[string]interface{}{"pubtoken_id": pubTokenIdString})
		}

		var pubTokenId uuid.UUID
		pubTokenId, err = uuid.Parse(pubTokenIdString)

		if err == nil {
			var pubToken *pubtokenmodel.PubTokenModel
			pubToken, err = service.pubTokenService.Detail(c.UserContext(), pubTokenId)

			if err == nil {
				c.Locals("read", pubToken.Read)
				c.Locals("write", pubToken.Write)
			}
		}
	}

	return err
}
func (service *pubTokenMiddlewareImpl) CanRefresh(c *fiber.Ctx) error {
	return fiber.NewError(401, "Unauthenticated")
}

// impl `jwt.JwtMiddleware` end
