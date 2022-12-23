package repositories

import (
	"context"
	"mime/multipart"
	"outclass-api/app/constants"
	"outclass-api/app/models"
	"outclass-api/app/repositories/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DirectoryRepositories struct {
	*mongo.Client
}

const storedFilePath = "./uploads/files/"

func (r *DirectoryRepositories) CreatePost(c *fiber.Ctx, directory models.Directory, fileHeaders []*multipart.FileHeader) ([]models.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fileModels := []models.File{}
	for _, fileHeader := range fileHeaders {
		fileModel, err := helpers.UploadFile(c, fileHeader, storedFilePath)
		if err != nil {
			for _, fileModel := range fileModels {
				helpers.DeleteFile(fileModel.Id, storedFilePath)
			}
			return nil, err
		}
		fileModels = append(fileModels, *fileModel)
	}

	directory.Files = append(directory.Files, fileModels...)
	_, err := helpers.DirectoryCollection(r.Client).InsertOne(ctx, directory)
	if err != nil {
		for _, fileModel := range fileModels {
			helpers.DeleteFile(fileModel.Id, storedFilePath)
		}
		return nil, err
	}
	return fileModels, nil
}

func (r *DirectoryRepositories) CreateFolder(directory models.Directory) error {
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
