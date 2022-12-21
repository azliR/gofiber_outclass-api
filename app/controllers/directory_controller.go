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
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	createPostDto := &dtos.CreatePostDto{}

	if err := c.BodyParser(createPostDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := validator.Struct(createPostDto); err != nil {
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

	if createPostDto.ParentId != "" {
		foundedDirectory, err := db.GetDirectoryById(createPostDto.ParentId)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return c.Status(fiber.StatusNotFound).JSON(commons.Response{
					Success: false,
					Message: "the parent document is not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		} else {
			if foundedDirectory.Type != constants.FolderDirectoryName {
				return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
					Success: false,
					Message: "the parent is not a folder",
				})
			}
		}
	}

	if createPostDto.ClassroomId != "" {
		_, err = db.GetClassroomById(createPostDto.ClassroomId)
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
	parentId, _ := primitive.ObjectIDFromHex(createPostDto.ParentId)
	ownerId, _ := primitive.ObjectIDFromHex(claims.UserId)
	classroomId, _ := primitive.ObjectIDFromHex(createPostDto.ClassroomId)
	directory := models.Directory{
		Id:           primitive.NewObjectID(),
		ParentId:     &parentId,
		OwnerId:      ownerId,
		ClassroomId:  &classroomId,
		Name:         createPostDto.Name,
		Type:         constants.PostDirectoryName,
		Description:  createPostDto.Description,
		Files:        []models.File{},
		LastModified: primitive.NewDateTimeFromTime(now),
		DateCreated:  primitive.NewDateTimeFromTime(now),
	}

	if err := validator.Struct(directory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
			Success: false,
			Message: utils.ValidatorErrors(err),
		})
	}

	err, files := db.CreatePost(c, directory, fileHeaders)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	directory.Files = append(directory.Files, files...)
	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    responses.ToPostResponse(directory),
	})
}

func CreateFolder(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
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
					Message: "the parent document is not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		} else {
			if foundedDirectory.Type != constants.FolderDirectoryName {
				return c.Status(fiber.StatusBadRequest).JSON(commons.Response{
					Success: false,
					Message: "the parent is not a folder",
				})
			}
		}
	}

	if folderDto.ClassroomId != "" {
		_, err = db.GetClassroomById(folderDto.ClassroomId)
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
	parentId, _ := primitive.ObjectIDFromHex(folderDto.ParentId)
	ownerId, _ := primitive.ObjectIDFromHex(claims.UserId)
	classroomId, _ := primitive.ObjectIDFromHex(folderDto.ClassroomId)
	color := folderDto.Color
	if color == "" {
		color = "grape"
	}
	directory := models.Directory{
		Id:           primitive.NewObjectID(),
		ParentId:     &parentId,
		OwnerId:      ownerId,
		ClassroomId:  &classroomId,
		Name:         folderDto.Name,
		Type:         constants.FolderDirectoryName,
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

	if err := db.CreateFolder(directory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(commons.Response{
		Success: true,
		Data:    responses.ToFolderResponse(directory),
	})
}

func GetDirectoryById(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
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

	var directoryResponse interface{}
	if directory.Type == constants.PostDirectoryName {
		directoryResponse = responses.ToPostResponse(*directory)
	} else {
		directoryResponse = responses.ToFolderResponse(*directory)
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    directoryResponse,
	})

}

func GetDirectoriesByParentId(c *fiber.Ctx) error {
	claims, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
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

	if directoryParam.ShareType == constants.ClassShareTypeName {
		if _, err := db.GetClassroomMemberByUserIdAndClassroomId(claims.UserId, directoryParam.ClassroomId); err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return c.Status(fiber.StatusNotFound).JSON(commons.Response{
					Success: false,
					Message: "you are not a member of this classroom",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: err.Error(),
			})
		}
	}

	directories, err := db.GetDirectoriesByParentId(
		claims.UserId,
		directoryParam.ClassroomId,
		directoryParam.Type,
		directoryParam.ShareType,
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

	responseDirectories := make([]interface{}, 0)
	for _, directory := range *directories {
		var directoryResponse interface{}
		if directory.Type == constants.PostDirectoryName {
			directoryResponse = responses.ToPostResponse(directory)
		} else {
			directoryResponse = responses.ToFolderResponse(directory)
		}
		responseDirectories = append(responseDirectories, directoryResponse)
	}

	return c.Status(fiber.StatusOK).JSON(commons.Response{
		Success: true,
		Data:    responseDirectories,
	})
}

func UpdatePostById(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
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
		if foundedDirectory.Type != constants.PostDirectoryName {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: "this directory is not a post",
			})
		}
	}

	now := time.Now()
	parentId, _ := primitive.ObjectIDFromHex(postDto.ParentId)

	foundedDirectory.ParentId = &parentId

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
		Data:    responses.ToPostResponse(*foundedDirectory),
	})
}

func UpdateFolderById(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
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
		if foundedDirectory.Type != constants.FolderDirectoryName {
			return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
				Success: false,
				Message: "this directory is not a folder",
			})
		}
	}

	now := time.Now()
	parentId, _ := primitive.ObjectIDFromHex(folderDto.ParentId)

	foundedDirectory.ParentId = &parentId
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
		Data:    responses.ToFolderResponse(*foundedDirectory),
	})
}

func DeleteDirectory(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
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

// func UploadFile(c *fiber.Ctx) error {
// 	claims, err := core.VerifyAndSyncToken(c)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}

// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 	}

// 	fileModels := []models.File{}
// 	for _, fileHeaders := range form.File {
// 		for _, fileHeader := range fileHeaders {
// 			fileModel, err := repositories.UploadFile(c, fileHeader, claims.UserId)
// 			if err != nil {
// 				return c.Status(fiber.StatusInternalServerError).JSON(commons.Response{
// 					Success: false,
// 					Message: err.Error(),
// 				})
// 			}
// 			fileModels = append(fileModels, *fileModel)
// 		}
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(commons.Response{
// 		Success: false,
// 		Data:    responses.ToFileResponses(fileModels),
// 	})
// }

func GetFile(c *fiber.Ctx) error {
	_, err := core.VerifyAndSyncToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
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
		return c.Status(fiber.StatusUnauthorized).JSON(commons.Response{
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
