package repositories

import (
	"context"
	"mime/multipart"
	"os"
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

func (r *DirectoryRepositories) GetDirectoriesByParentId(parentId string, page uint, pageLimit uint) (*[]models.Directory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(parentId)

	findOption := options.Find()
	findOption.SetSort(bson.D{
		{Key: "type", Value: 1},
		{Key: "last_modified", Value: -1},
	})
	findOption.SetLimit(int64(pageLimit))
	findOption.SetSkip(int64((page - 1) * pageLimit))

	cursor, err := helpers.DirectoryCollection(r.Client).Find(ctx, bson.M{"_parent_id": objId}, findOption)
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
	if err := c.SaveFile(file, filePath); err != nil {
		return nil, err
	}
	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")
	server := "http://" + serverHost + ":" + serverPort
	link := server + "/api/v1/files/" + fileName

	return &models.File{
		Link: link,
		Type: strings.Replace(fileExt, ".", "", 1),
		Size: file.Size,
	}, nil
}

func DeleteFile(c *fiber.Ctx, fileId string) error {
	filePath := StoredFilePath + fileId

	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}
