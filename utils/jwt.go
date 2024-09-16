package utils

import (
	jwtservice "gofiber-boilerplate/modules/jwt"

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

func GetFiberJwtUserId(c *fiber.Ctx) (id uuid.UUID, err error) {
	idString, err := GetFiberJwtUserIdString(c)
	if err != nil {
		return
	}

	id, err = uuid.Parse(idString)
	return
}
