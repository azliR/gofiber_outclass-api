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

func CreateEvent(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	eventDto := &dtos.CreateEventDto{}
	print(eventDto.Repeat)

	if err := c.BodyParser(eventDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(eventDto); err != nil {
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

	if eventDto.ClassroomId != "" {
		_, err = db.GetClassroomById(eventDto.ClassroomId)
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
	classroomId, _ := primitive.ObjectIDFromHex(eventDto.ClassroomId)

	var endDate *primitive.DateTime
	if eventDto.EndDate != nil {
		newEndDate := primitive.NewDateTimeFromTime(*eventDto.EndDate)
		endDate = &newEndDate
	}
	color := eventDto.Color
	if color == "" {
		color = "grape"
	}
	event := models.Event{
		Id:           primitive.NewObjectID(),
		CreatorId:    creatorId,
		ClassroomId:  &classroomId,
		Name:         eventDto.Name,
		StartDate:    primitive.NewDateTimeFromTime(eventDto.StartDate),
		EndDate:      endDate,
		Repeat:       eventDto.Repeat,
		Color:        color,
		Description:  eventDto.Description,
		LastModified: primitive.NewDateTimeFromTime(now),
		DateCreated:  primitive.NewDateTimeFromTime(now),
	}

	if err := validator.Struct(event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	if err := db.CreateEvent(event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    responses.ToEventResponse(event),
	})
}

func GetEventsByClassroomId(c *fiber.Ctx) error {
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

	events, err := db.GetEventsByClassroomId(classroomId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    responses.ToEventResponses(*events),
	})
}
