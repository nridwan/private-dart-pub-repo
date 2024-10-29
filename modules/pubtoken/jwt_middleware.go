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
	CanWrite(c *fiber.Ctx) error
	GetPubUserId(c *fiber.Ctx) uuid.UUID
	HasAccess(c *fiber.Ctx) error
}

type pubTokenMiddlewareImpl struct {
	jwtService      jwt.JwtService
	pubTokenService PubTokenService
	monitorService  monitor.MonitorService
}

func NewPubTokenJwtMiddleware(jwtService jwt.JwtService, pubTokenService PubTokenService, monitorService monitor.MonitorService) PubTokenJwtMiddleware {
	return &pubTokenMiddlewareImpl{
		jwtService:      jwtService,
		pubTokenService: pubTokenService,
		monitorService:  monitorService,
	}
}

// impl `PubTokenJwtMiddleware` start

func (service *pubTokenMiddlewareImpl) CanWrite(c *fiber.Ctx) error {
	if c.Locals("write") != true {
		return fiber.NewError(401, "Unauthenticated")
	}

	return c.Next()
}

func (service *pubTokenMiddlewareImpl) GetPubUserId(c *fiber.Ctx) uuid.UUID {
	return c.Locals("pub_user_id").(uuid.UUID)
}

// impl `PubTokenJwtMiddleware` end

func (service *pubTokenMiddlewareImpl) HasAccess(c *fiber.Ctx) error {
	err := service.jwtService.CanAccess(c, JwtIssuer)

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
				c.Locals("write", pubToken.Write)
				c.Locals("pub_user_id", *pubToken.UserID)
			}
		}
	}

	return err
}

// impl `jwt.JwtMiddleware` start

func (service *pubTokenMiddlewareImpl) CanAccess(c *fiber.Ctx) error {
	err := service.HasAccess(c)

	if err == nil {
		return c.Next()
	}

	return err
}

func (service *pubTokenMiddlewareImpl) CanRefresh(c *fiber.Ctx) error {
	return fiber.NewError(401, "Unauthenticated")
}

// impl `jwt.JwtMiddleware` end
