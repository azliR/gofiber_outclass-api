package core

import (
	"outclass-api/app/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func VerifyAndSyncToken(c *fiber.Ctx) (*utils.TokenMetadata, error) {
	now := time.Now().Unix()
	bearToken := c.Get("Authorization")
	token := utils.ExtractToken(bearToken)

	claims, err := utils.ExtractTokenMetadata(token)
	if err != nil {
		return nil, err
	}
	expires := claims.Expires

	if now > expires {
		return nil, errors.New("unauthorized, check expiration time of your token")
	}

	return claims, nil
}
