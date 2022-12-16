package controllers

import (
	"errors"
	"outclass-api/app/commons"
	"outclass-api/app/configs"
	"outclass-api/app/constants"
	"outclass-api/app/controllers/core"
	"outclass-api/app/controllers/responses"
	"outclass-api/app/dtos"
	"outclass-api/app/models"
	"outclass-api/app/utils"

	"github.com/dchest/uniuri"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateClassroom(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	classroomDto := &dtos.CreateClassroomDto{}

	if err := c.BodyParser(classroomDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(classroomDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	classroom := models.Classroom{
		Id:          primitive.NewObjectID(),
		Name:        classroomDto.Name,
		ClassCode:   uniuri.New(),
		Description: classroomDto.Description,
	}

	if err := validator.Struct(classroom); err != nil {
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

	userId, err := primitive.ObjectIDFromHex(claims.UserId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	foundedUser, err := db.GetUserById(claims.UserId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(commons.Response{
				Success: false,
				Message: "the user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	role, err := dtos.ToModelRole(constants.OwnerClassroomRoleName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	classroomMember := models.ClassroomMember{
		Id:            primitive.NewObjectID(),
		UserId:        userId,
		ClassroomId:   classroom.Id,
		StudentId:     classroomDto.StudentId,
		ClassroomName: classroom.Name,
		Name:          foundedUser.Name,
		Role:          role,
	}

	if err := validator.Struct(classroomMember); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	if err := db.CreateClassroom(classroom, classroomMember); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    responses.ToClassroomResponse(classroom),
	})
}

func GetClassroomById(c *fiber.Ctx) error {
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

	classroom, err := db.GetClassroomById(classroomId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    responses.ToClassroomResponse(*classroom),
	})
}

func GetClassroomByClassCode(c *fiber.Ctx) error {
	classCode := c.Params("classCode")

	db, err := configs.GetMongoConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	classroom, err := db.GetClassroomByClassCode(classCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    responses.ToClassroomResponse(*classroom),
	})
}

func JoinClassroomByClassCode(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	joinClassroomDto := &dtos.JoinClassroomDto{}

	if err := c.BodyParser(joinClassroomDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(joinClassroomDto); err != nil {
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

	foundedUser, err := db.GetUserById(claims.UserId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(commons.Response{
				Success: false,
				Message: "the user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	foundedClassroom, err := db.GetClassroomByClassCode(joinClassroomDto.ClassCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	userId, err := primitive.ObjectIDFromHex(claims.UserId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	role, err := dtos.ToModelRole(constants.MemberClassroomRoleName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	classroomMember := models.ClassroomMember{
		Id:            primitive.NewObjectID(),
		UserId:        userId,
		ClassroomId:   foundedClassroom.Id,
		StudentId:     joinClassroomDto.StudentId,
		ClassroomName: foundedClassroom.Name,
		Name:          foundedUser.Name,
		Role:          role,
	}
	if err := validator.Struct(classroomMember); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	if err = db.CreateClassroomMember(classroomMember); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
				Success: false,
				Message: "either you are already a student of this class or your student id is already taken by other student in this class",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    responses.ToClassroomMemberResponse(classroomMember),
	})
}

func UpdateClassroomById(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	classroomId := c.Params("classroomId")
	classroomDto := &dtos.UpdateClassroomDto{}
	if err := c.BodyParser(classroomDto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(classroomDto); err != nil {
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

	foundedClassroom, err := db.GetClassroomById(classroomId)
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

	if classroomDto.Name == foundedClassroom.Name && classroomDto.Description == foundedClassroom.Description {
		return c.SendStatus(fiber.StatusNotModified)
	} else {
		if classroomDto.Name != "" {
			foundedClassroom.Name = classroomDto.Name
		}
		foundedClassroom.Description = classroomDto.Description
	}

	if err = db.UpdateClassroom(*foundedClassroom); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    responses.ToClassroomResponse(*foundedClassroom),
	})
}

func DeleteClassroom(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	classroomId := c.Params("classroomId")
	_, err = primitive.ObjectIDFromHex(classroomId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	db, err := configs.GetMongoConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err = db.DeleteClassroom(classroomId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
