package controllers

import (
	"context"
	"os"
	"outclass-api/app/commons"
	"outclass-api/app/configs"
	_auth "outclass-api/app/controllers/auth"
	"outclass-api/app/controllers/core"
	"outclass-api/app/dtos"
	"outclass-api/app/models"
	"outclass-api/app/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validator = utils.Validator

func UserSignUp(c *fiber.Ctx) error {
	signUpDto := &dtos.SignUp{}

	if err := c.BodyParser(signUpDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(signUpDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	user := &models.User{
		Id:       primitive.NewObjectID(),
		Name:     signUpDto.Name,
		Email:    signUpDto.Email,
		Password: utils.GeneratePassword(signUpDto.Password),
	}

	if err := validator.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	db, err := configs.GetMongoConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	foundedUser, err := db.GetUserByEmail(user.Email)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}
	} else {
		if foundedUser != nil {
			return c.Status(fiber.StatusConflict).JSON(commons.Response{
				Success: false,
				Message: "duplicate email",
			})
		}
	}

	if err = db.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    _auth.ToUserResponse(user),
	})
}

func UserSignIn(c *fiber.Ctx) error {
	signInDto := &dtos.SignIn{}

	if err := c.BodyParser(signInDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(signInDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	db, err := configs.GetMongoConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	foundedUser, err := db.GetUserByEmail(signInDto.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(commons.Response{
			Success: false,
			Message: "user with the given email is not found",
		})
	}

	compareUserPassword := utils.ComparePasswords(foundedUser.Password, signInDto.Password)
	if !compareUserPassword {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: "email and password do not match",
		})
	}

	userId := foundedUser.Id.Hex()
	tokens, err := utils.GenerateNewTokens(userId)
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
		Success: true,
		Data:    _auth.ToTokenResponse(tokens),
	})
}

func UserProfile(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	userId := claims.UserId

	db, err := configs.GetMongoConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	user, err := db.GetUserById(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    _auth.ToUserResponse(user),
	})
}

func UserSignOut(c *fiber.Ctx) error {
	bearToken := c.Get("Authorization")
	token := utils.ExtractToken(bearToken)

	claims, err := utils.ExtractTokenMetadata(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	userID := claims.UserId

	redis, err := configs.GetRedisConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	err = redis.Del(context.Background(), userID).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
