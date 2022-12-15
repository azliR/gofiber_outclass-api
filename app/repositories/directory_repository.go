package repositories

import (
	"context"
	"errors"
	"mime/multipart"
	"os"
	"outclass-api/app/constants"
	"outclass-api/app/models"
	"outclass-api/app/repositories/helpers"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DirectoryRepositories struct {
	*mongo.Client
}

const StoredFilePath = "./uploads/"

func (r *DirectoryRepositories) CreateDirectory(directory models.Directory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := helpers.DirectoryCollection(r.Client).InsertOne(ctx, directory)
	if err != nil {
		return err
	}
	return nil
}

func (r *DirectoryRepositories) GetDirectoryById(directoryId string) (*models.Directory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	directory := &models.Directory{}

	objId, _ := primitive.ObjectIDFromHex(directoryId)

	err := helpers.DirectoryCollection(r.Client).FindOne(ctx, bson.M{"_id": objId}).Decode(directory)
	if err != nil {
		return nil, err
	}
	return directory, nil
}

func (r *DirectoryRepositories) GetDirectoriesByParentId(
	ownerId string,
	classroomId string,
	directoryType string,
	shareType string,
	parentId string,
	page uint,
	pageLimit uint,
) (*[]models.Directory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parentObjId, _ := primitive.ObjectIDFromHex(parentId)
	sortKey := "name"
	sortValue := 1
	if directoryType == constants.PostDirectoryName {
		sortKey = "created_at"
		sortValue = -1
	}

	findOption := options.Find()
	findOption.SetSort(bson.D{
		{Key: sortKey, Value: sortValue},
	})
	findOption.SetLimit(int64(pageLimit))
	findOption.SetSkip(int64((page - 1) * pageLimit))

	shareTypeFilter := bson.E{}
	if shareType == constants.ClassShareTypeName {
		classroomIdObj, err := primitive.ObjectIDFromHex(classroomId)
		if err != nil {
			return nil, err
		}
		shareTypeFilter = bson.E{
			Key:   "_classroom_id",
			Value: bson.M{"$eq": classroomIdObj},
		}
	} else {
		ownerIdObj, err := primitive.ObjectIDFromHex(ownerId)
		if err != nil {
			return nil, err
		}

		if shareType == constants.GroupShareTypeName {
			shareTypeFilter = bson.E{
				Key: "$and",
				Value: bson.A{
					bson.M{
						"$or": bson.A{
							bson.M{"_owner_id": ownerIdObj},
							bson.M{"shared_with._user_id": ownerIdObj},
						},
					},
					bson.M{"shared_with": bson.M{"$ne": nil}},
					bson.M{"_classroom_id": primitive.NilObjectID},
				},
			}
		} else if shareType == constants.PersonalShareTypeName {
			shareTypeFilter = bson.E{
				Key: "$and",
				Value: bson.A{
					bson.M{"_owner_id": ownerIdObj},
					bson.M{"_classroom_id": primitive.NilObjectID},
					bson.M{"shared_with": nil},
				},
			}
		}
	}

	cursor, err := helpers.DirectoryCollection(r.Client).Find(
		ctx,
		bson.D{
			{Key: "_parent_id", Value: parentObjId},
			{Key: "type", Value: directoryType},
			shareTypeFilter,
		},
		findOption,
	)
	if err != nil {
		return nil, err
	}

	directories := &[]models.Directory{}
	cursor.All(ctx, directories)

	return directories, nil
}

func (r *DirectoryRepositories) UpdateDirectory(directory models.Directory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := helpers.DirectoryCollection(r.Client).UpdateByID(ctx, directory.Id, bson.D{
		{Key: "$set", Value: directory},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *DirectoryRepositories) DeleteDirectory(directoryId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(directoryId)
	_, err := helpers.DirectoryCollection(r.Client).DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return err
	}
	return nil
}

func UploadFile(c *fiber.Ctx, file *multipart.FileHeader, userId string) (*models.File, error) {
	fileExt := filepath.Ext(file.Filename)
	fileName := userId + strings.ReplaceAll(uuid.NewString(), "-", "") + fileExt
	filePath := StoredFilePath + fileName
	if _, err := os.Stat(StoredFilePath); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(StoredFilePath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	if err := c.SaveFile(file, filePath); err != nil {
		return nil, err
	}
	serverScheme := os.Getenv("SERVER_SCHEME")
	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")
	server := serverScheme + "://" + serverHost + ":" + serverPort
	link := server + "/api/v1/files/" + fileName
	fileType := strings.Replace(fileExt, ".", "", 1)

	return &models.File{
		Name: file.Filename,
		Link: link,
		Type: &fileType,
		Size: &file.Size,
	}, nil
}

func DeleteFile(c *fiber.Ctx, fileId string) error {
	filePath := StoredFilePath + fileId

	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}
