package controllers

import (
	"errors"
	"outclass-api/app/commons"
	"outclass-api/app/configs"
	"outclass-api/app/controllers/core"
	_directory "outclass-api/app/controllers/directory"
	"outclass-api/app/dtos"
	"outclass-api/app/models"
	"outclass-api/app/repositories"
	"outclass-api/app/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePost(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	postDto := &dtos.CreatePostDto{}

	if err := c.BodyParser(postDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(postDto); err != nil {
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

	if postDto.ParentId != "" {
		foundedDirectory, err := db.GetDirectoryById(postDto.ParentId)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return c.Status(fiber.StatusNotFound).JSON(commons.Response{
					Success: false,
					Message: "The parent document is not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		} else {
			if foundedDirectory.Type != "folder" {
				return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
					Success: false,
					Message: "The parent is not a folder",
				})
			}
		}
	}
	// TODO: implement checking classroomId

	now := time.Now()
	parentId, _ := primitive.ObjectIDFromHex(postDto.ParentId)
	ownerId, _ := primitive.ObjectIDFromHex(claims.UserId)
	classroomId, _ := primitive.ObjectIDFromHex(postDto.ClassroomId)
	directory := models.Directory{
		Id:           primitive.NewObjectID(),
		ParentId:     parentId,
		OwnerId:      ownerId,
		ClassroomId:  classroomId,
		Name:         postDto.Name,
		Type:         "post",
		Description:  postDto.Description,
		Files:        dtos.ToModelFiles(postDto.Files),
		LastModified: primitive.NewDateTimeFromTime(now),
		DateCreated:  primitive.NewDateTimeFromTime(now),
	}

	if err := validator.Struct(directory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	if err := db.CreateDirectory(directory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    directory,
	})
}

func CreateFolder(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	folderDto := &dtos.CreateFolderDto{}

	if err := c.BodyParser(folderDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(folderDto); err != nil {
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

	if folderDto.ParentId != "" {
		foundedDirectory, err := db.GetDirectoryById(folderDto.ParentId)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return c.Status(fiber.StatusNotFound).JSON(commons.Response{
					Success: false,
					Message: "The parent document is not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		} else {
			if foundedDirectory.Type != "folder" {
				return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
					Success: false,
					Message: "The parent is not a folder",
				})
			}
		}
	}
	// TODO: implement checking classroomId

	now := time.Now()
	parentId, _ := primitive.ObjectIDFromHex(folderDto.ParentId)
	ownerId, _ := primitive.ObjectIDFromHex(claims.UserId)
	classroomId, _ := primitive.ObjectIDFromHex(folderDto.ClassroomId)
	color := folderDto.Color
	if color == "" {
		color = "teal"
	}
	directory := models.Directory{
		Id:           primitive.NewObjectID(),
		ParentId:     parentId,
		OwnerId:      ownerId,
		ClassroomId:  classroomId,
		Name:         folderDto.Name,
		Type:         "folder",
		Color:        &color,
		Description:  folderDto.Description,
		LastModified: primitive.NewDateTimeFromTime(now),
		DateCreated:  primitive.NewDateTimeFromTime(now),
	}

	if err := validator.Struct(directory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	if err := db.CreateDirectory(directory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    directory,
	})
}

func GetDirectoryById(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	directoryId := c.Params("directoryId")

	db, err := configs.GetMongoConnection()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	directory, err := db.GetDirectoryById(directoryId)
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
		Data:    directory,
	})
}

func GetDirectoriesByParentId(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	directoryParam := &dtos.GetDirectoriesDto{}

	if err = c.QueryParser(directoryParam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(directoryParam); err != nil {
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

	directories, err := db.GetDirectoriesByParentId(
		directoryParam.ParentId,
		directoryParam.Page,
		directoryParam.PageLimit,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    directories,
	})
}

func UpdatePostById(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	postId := c.Params("postId")
	postDto := &dtos.UpdatePostDto{}
	if err := c.BodyParser(postDto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(postDto); err != nil {
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

	foundedDirectory, err := db.GetDirectoryById(postId)
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
	} else {
		if foundedDirectory.Type != "post" {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: "This directory is not a post",
			})
		}
	}

	now := time.Now()
	parentId, _ := primitive.ObjectIDFromHex(postDto.ParentId)

	foundedDirectory.ParentId = parentId

	if postDto.Name == foundedDirectory.Name &&
		postDto.Description == foundedDirectory.Description &&
		postDto.ParentId == foundedDirectory.ParentId.Hex() {
		return c.SendStatus(fiber.StatusNotModified)
	} else {
		if postDto.Name != "" {
			foundedDirectory.Name = postDto.Name
		}
		foundedDirectory.Description = postDto.Description
		foundedDirectory.LastModified = primitive.NewDateTimeFromTime(now)
	}

	if err = db.UpdateDirectory(*foundedDirectory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    foundedDirectory,
	})
}

func UpdateFolderById(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	folderId := c.Params("folderId")
	folderDto := &dtos.UpdateFolderDto{}
	if err := c.BodyParser(folderDto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(folderDto); err != nil {
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

	foundedDirectory, err := db.GetDirectoryById(folderId)
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
	} else {
		if foundedDirectory.Type != "folder" {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: "This directory is not a folder",
			})
		}
	}

	now := time.Now()
	parentId, _ := primitive.ObjectIDFromHex(folderDto.ParentId)

	foundedDirectory.ParentId = parentId
	if folderDto.Name == foundedDirectory.Name &&
		folderDto.Color == foundedDirectory.Color &&
		folderDto.Description == foundedDirectory.Description &&
		folderDto.ParentId == foundedDirectory.ParentId.Hex() {
		return c.SendStatus(fiber.StatusNotModified)
	} else {
		if folderDto.Name != "" {
			foundedDirectory.Name = folderDto.Name
		}
		if folderDto.Color != nil {
			foundedDirectory.Color = folderDto.Color
		}
		foundedDirectory.Description = folderDto.Description
		foundedDirectory.LastModified = primitive.NewDateTimeFromTime(now)
	}

	if err = db.UpdateDirectory(*foundedDirectory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    foundedDirectory,
	})
}

func DeleteDirectory(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	directoryId := c.Params("directoryId")
	_, err = primitive.ObjectIDFromHex(directoryId)
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

	if err = db.DeleteDirectory(directoryId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func UploadFile(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	fileModel, err := repositories.UploadFile(c, file, claims.UserId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: false,
		Data:    _directory.ToFileUploadResponse(*fileModel),
	})
}

func GetFile(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	fileId := c.Params("fileId")

	return c.SendFile("./uploads/"+fileId, true)
}

func DeleteFile(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	fileId := c.Params("fileId")

	if err := repositories.DeleteFile(c, fileId); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
