package controllers

import (
	"errors"
	"outclass-api/app/commons"
	"outclass-api/app/configs"
	"outclass-api/app/controllers/core"
	"outclass-api/app/controllers/responses"
	"outclass-api/app/dtos"
	"outclass-api/app/models"
	"outclass-api/app/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTask(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	taskDto := &dtos.CreateTaskDto{}

	if err := c.BodyParser(taskDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(taskDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	forms, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	fileHeaders := forms.File["files"]

	db, err := configs.GetMongoConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if taskDto.ClassroomId != "" {
		_, err = db.GetClassroomById(taskDto.ClassroomId)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
					Success: false,
					Message: "the classroomId is not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}
	}

	now := time.Now()
	creatorId, _ := primitive.ObjectIDFromHex(claims.UserId)
	classroomId, _ := primitive.ObjectIDFromHex(taskDto.ClassroomId)

	color := taskDto.Color
	if color == "" {
		color = "grape"
	}
	task := models.Task{
		Id:           primitive.NewObjectID(),
		CreatorId:    creatorId,
		ClassroomId:  &classroomId,
		Title:        taskDto.Title,
		Details:      taskDto.Details,
		Date:         primitive.NewDateTimeFromTime(taskDto.Date),
		Repeat:       taskDto.Repeat,
		Color:        color,
		LastModified: primitive.NewDateTimeFromTime(now),
		DateCreated:  primitive.NewDateTimeFromTime(now),
	}

	if err := validator.Struct(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	files, err := db.CreateTask(c, task, fileHeaders)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	task.Files = append(task.Files, files...)
	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    responses.ToTaskResponse(task),
	})
}

func GetTasksByClassroomId(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	classroomId := c.Params("classroomId")

	db, err := configs.GetMongoConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	tasks, err := db.GetTasksByClassroomId(classroomId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    responses.ToTaskResponses(*tasks),
	})
}
