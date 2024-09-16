package jwt

import (
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type jwtCommonMiddleware interface {
	CanAccess(c *fiber.Ctx, issuer string) error
	CanRefresh(c *fiber.Ctx, issuer string) error
}

type JwtMiddleware interface {
	CanAccess(c *fiber.Ctx) error
	CanRefresh(c *fiber.Ctx) error
}

// impl `jwtCommonMiddleware` start

func (service *JwtModule) checkUser(c *fiber.Ctx, issuer string, refresh bool) bool {
	claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)

	var audiences jwt.ClaimStrings
	var err error

	if audiences, err = claims.GetAudience(); err != nil {
		return false
	}

	if iss, err := claims.GetIssuer(); err != nil || iss != issuer {
		return false
	}

	if (refresh && !slices.Contains(audiences, JwtRefreshAud)) || (!refresh && !slices.Contains(audiences, JwtAppAud)) {
		return false
	}

	if _, err = claims.GetSubject(); err != nil {
		return false
	}

	return true
}

func (service *JwtModule) CanAccess(c *fiber.Ctx, issuer string) error {
	if service.checkUser(c, issuer, false) {
		return c.Next()
	}
	return fiber.NewError(401, "Unauthenticated")
}

func (service *JwtModule) CanRefresh(c *fiber.Ctx, issuer string) error {
	if service.checkUser(c, issuer, true) {
		return c.Next()
	}
	return fiber.NewError(401, "Unauthenticated")
}

// impl `jwtCommonMiddleware` end
