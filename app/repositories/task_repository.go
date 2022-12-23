package repositories

import (
	"context"
	"mime/multipart"
	"outclass-api/app/models"
	"outclass-api/app/repositories/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepositories struct {
	*mongo.Client
}

const storedTaskFilePath = "./uploads/tasks/"

func (r *TaskRepositories) CreateTask(c *fiber.Ctx, task models.Task, fileHeaders []*multipart.FileHeader) ([]models.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fileModels := []models.File{}
	for _, fileHeader := range fileHeaders {
		fileModel, err := helpers.UploadFile(c, fileHeader, storedTaskFilePath)
		if err != nil {
			for _, fileModel := range fileModels {
				helpers.DeleteFile(fileModel.Id, storedTaskFilePath)
			}
			return nil, err
		}
		fileModels = append(fileModels, *fileModel)
	}

	task.Files = append(task.Files, fileModels...)
	_, err := helpers.TaskCollection(r.Client).InsertOne(ctx, task)
	if err != nil {
		for _, fileModel := range fileModels {
			helpers.DeleteFile(fileModel.Id, storedTaskFilePath)
		}
		return nil, err
	}
	return fileModels, nil
}

func (r *TaskRepositories) GetTaskById(taskId string) (*models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	task := &models.Task{}

	objId, _ := primitive.ObjectIDFromHex(taskId)

	err := helpers.TaskCollection(r.Client).FindOne(ctx, bson.M{"_id": objId}).Decode(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *TaskRepositories) GetTasksByClassroomId(classroomId string) (*[]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(classroomId)

	cursor, err := helpers.TaskCollection(r.Client).Find(ctx, bson.M{"_classroom_id": objId})
	if err != nil {
		return nil, err
	}

	tasks := &[]models.Task{}
	cursor.All(ctx, tasks)

	return tasks, nil
}

func (r *TaskRepositories) UpdateTask(task models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := helpers.TaskCollection(r.Client).UpdateByID(ctx, task.Id, bson.D{
		{Key: "$set", Value: task},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepositories) DeleteTask(taskId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(taskId)
	_, err := helpers.TaskCollection(r.Client).DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return err
	}
	return nil
}
