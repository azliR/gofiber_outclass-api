package controllers

import (
	"context"
	"os"
	"outclass-api/app/commons"
	"outclass-api/app/configs"
	"outclass-api/app/controllers/responses"
	"outclass-api/app/dtos"
	"outclass-api/app/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RenewTokens(c *fiber.Ctx) error {
	now := time.Now().Unix()
	bearToken := c.Get("Authorization")
	token := utils.ExtractToken(bearToken)

	claims, err := utils.ExtractTokenMetadata(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	expiresAccessToken := claims.Expires

	if now < expiresAccessToken {
		return c.Status(fiber.StatusForbidden).JSON(commons.Response{
			Success: false,
			Message: "your token still active",
		})
	}

	renew := &dtos.Renew{}

	if err := c.BodyParser(renew); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	expiresRefreshToken, err := utils.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if now < expiresRefreshToken {
		userId := claims.UserId

		tokens, err := utils.GenerateNewTokens(userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		redis, err := configs.GetRedisConnection()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		// renewToken, err := redis.Get(context.Background(), userId).Result()
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
		// 		Success: false,
		// 		Message: err.Error(),
		// 	})
		// } else {
		// 	if renewToken != renew.RefreshToken {
		// 		return c.Status(fiber.StatusForbidden).JSON(commons.Response{
		// 			Success: false,
		// 			Message: "token already used",
		// 		})
		// 	}
		// }

		hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))
		expiration := time.Hour * time.Duration(hoursCount)

		err = redis.Set(context.Background(), userId, tokens.Refresh, expiration).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(commons.Response{
			Success: true,
			Data: responses.TokenResponse{
				AccessToken:           tokens.Access,
				TokenExpiresIn:        tokens.TokenExpiresIn,
				RefreshToken:          tokens.Refresh,
				RefreshTokenExpiresIn: tokens.RefreshTokenExpiresIn,
			},
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: "unauthorized, your session was ended earlier",
		})
	}
}
