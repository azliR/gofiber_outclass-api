package controllers

import (
	"errors"
	"outclass-api/app/commons"
	"outclass-api/app/configs"
	"outclass-api/app/controllers/core"
	"outclass-api/app/controllers/responses"
	"outclass-api/app/dtos"
	"outclass-api/app/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// func CreateClassroomMember(c *fiber.Ctx) error {
// 	_, err := core.VerifyAndSyncToken(c)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}

// 	classroomMemberDto := &dtos.CreateClassroomMemberDto{}

// 	if err := c.BodyParser(classroomMemberDto); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}

// 	if err := validator.Struct(classroomMemberDto); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
// 			Success: false,
// 			Message: utils.ValidatorErrors(err),
// 		})
// 	}

// 	classroomIdStr := c.Params("classroomId")
// 	classroomId, err := primitive.ObjectIDFromHex(classroomIdStr)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}
// 	userId, err := primitive.ObjectIDFromHex(classroomMemberDto.UserId)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}
// 	role, err := dtos.ToModelRole(classroomMemberDto.Role)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}

// 	db, err := configs.GetMongoConnection()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}

// 	foundedClassroom, err := db.GetClassroomById(classroomIdStr)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return c.Status(fiber.StatusNotFound).JSON(commons.Response{
// 				Success: false,
// 				Message: "the classroom is not found",
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}

// 	foundedUser, err := db.GetUserById(classroomMemberDto.UserId)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return c.Status(fiber.StatusNotFound).JSON(commons.Response{
// 				Success: false,
// 				Message: "the user not found",
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}
// 	classroomMember := models.ClassroomMember{
// 		Id:            primitive.NewObjectID(),
// 		UserId:        userId,
// 		ClassroomId:   classroomId,
// 		StudentId:     classroomMemberDto.StudentId,
// 		Name:          foundedUser.Name,
// 		ClassroomName: foundedClassroom.Name,
// 		Role:          role,
// 	}

// 	if err := validator.Struct(classroomMember); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
// 			Success: false,
// 			Message: utils.ValidatorErrors(err),
// 		})
// 	}

// 	if err := db.CreateClassroomMember(classroomMember); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(commons.Response{
// 		Success: true,
// 		Data:    responses.ToClassroomMemberResponse(classroomMember),
// 	})
// }

func GetClassroomMemberByUserIdAndClassroomId(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)
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

	classroom, err := db.GetClassroomById(classroomId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(commons.Response{
				Success: false,
				Message: "the classroom is not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	classMembersCount, err := db.CountClassroomMembersByClassroomId(classroomId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(commons.Response{
				Success: false,
				Message: "the classroom is not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	classroomMember, err := db.GetClassroomMemberByUserIdAndClassroomId(claims.UserId, classroomId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(commons.Response{
				Success: false,
				Message: "the classroom member is not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    responses.ToClassroomMemberProfileResponse(*classroom, *classroomMember, classMembersCount),
	})
}

func GetClassroomMembersByClassroomId(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	classroomId := c.Params("classroomId")
	classroomMemberParam := &dtos.GetClassroomMembersDto{}

	if err = c.QueryParser(classroomMemberParam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(classroomMemberParam); err != nil {
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

	classroomMembers, err := db.GetClassroomMembersByClassroomId(
		classroomId,
		classroomMemberParam.Page,
		classroomMemberParam.PageLimit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    responses.ToClassroomMemberResponses(*classroomMembers),
	})
}

func GetClassroomMembersByUserId(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	classroomMemberParam := &dtos.GetClassroomMembersDto{}

	if err = c.QueryParser(classroomMemberParam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(classroomMemberParam); err != nil {
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

	classroomMembers, err := db.GetClassroomMembersByUserId(
		claims.UserId,
		classroomMemberParam.Page,
		classroomMemberParam.PageLimit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    responses.ToClassroomMemberResponses(*classroomMembers),
	})
}
