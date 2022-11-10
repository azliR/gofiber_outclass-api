package handlers

import (
	"context"
	"os"
	"outclass-api/app/commons"
	"outclass-api/app/configs"
	"outclass-api/app/dtos"
	_auth "outclass-api/app/handlers/auth"
	"outclass-api/app/handlers/core"
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
		Id:        primitive.NewObjectID(),
		StudentId: signUpDto.StudentId,
		Name:      signUpDto.Name,
		Email:     signUpDto.Email,
		Password:  utils.GeneratePassword(signUpDto.Password),
	}

	if err := validator.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	db := configs.MongoDb

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
				Message: "Duplicate email",
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

	db := configs.MongoDb

	foundedUser, err := db.GetUserByEmail(signInDto.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(commons.Response{
			Success: false,
			Message: "User with the given email is not found",
		})
	}

	compareUserPassword := utils.ComparePasswords(foundedUser.Password, signInDto.Password)
	if !compareUserPassword {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: "Email and password do not match",
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

	connRedis := configs.RedisDb
	errSaveToRedis := connRedis.Set(context.Background(), userId, tokens.Refresh, expiration).Err()

	if errSaveToRedis != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: errSaveToRedis.Error(),
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

	db := configs.MongoDb

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

	connRedis := configs.RedisDb

	errDelFromRedis := connRedis.Del(context.Background(), userID).Err()
	if errDelFromRedis != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: errDelFromRedis.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
