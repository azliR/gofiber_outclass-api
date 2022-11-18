package controllers

import (
	"context"
	"fmt"
	"os"
	"outclass-api/app/commons"
	"outclass-api/app/configs"
	_auth "outclass-api/app/controllers/auth"
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

	fmt.Print("Time Now: ")
	fmt.Println(time.Unix(now, 0))
	fmt.Print("Time Exp: ")
	fmt.Println(time.Unix(expiresAccessToken, 0))

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

		hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))
		expiration := time.Hour * time.Duration(hoursCount)

		redis, err := configs.GetRedisConnection()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		err = redis.Set(context.Background(), userId, tokens.Refresh, expiration).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(commons.Response{
			Success: false,
			Data: _auth.TokenResponse{
				AccessToken:  tokens.Access,
				RefreshToken: tokens.Refresh,
			},
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: "unauthorized, your session was ended earlier",
		})
	}
}
