package utils

import (
	jwtservice "private-pub-repo/modules/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetFiberJwtClaims(c *fiber.Ctx) jwtservice.JwtClaim {
	return c.Locals("user").(*jwt.Token).Claims.(jwtservice.JwtClaim)
}

func GetFiberJwtUserIdString(c *fiber.Ctx) (id string, err error) {
	return GetFiberJwtClaims(c).GetSubject()
}

func HasJwt(c *fiber.Ctx) bool {
	_, ok := c.Locals("user").(*jwt.Token)
	return ok
}

func IsFiberJwtCanWrite(c *fiber.Ctx) bool {
	raw, ok := GetFiberJwtClaims(c)["can_write"]

	if !ok {
		return false
	}

	result, ok := raw.(bool)

	return ok && result
}

func GetFiberJwtUserId(c *fiber.Ctx) (id uuid.UUID, err error) {
	idString, err := GetFiberJwtUserIdString(c)
	if err != nil {
		return
	}

	id, err = uuid.Parse(idString)
	return
}
