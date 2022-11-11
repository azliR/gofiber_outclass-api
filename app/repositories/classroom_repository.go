package repositories

import (
	"context"
	"outclass-api/app/models"
	"outclass-api/app/repositories/helpers"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassroomRepositories struct {
	*mongo.Client
}

func (r *ClassroomRepositories) CreateClassroom(classroom models.Classroom) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := helpers.ClassroomCollection(r.Client).InsertOne(ctx, classroom)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClassroomRepositories) GetClassroomById(classroomId string) (*models.Classroom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	classroom := &models.Classroom{}

	objId, _ := primitive.ObjectIDFromHex(classroomId)

	err := helpers.ClassroomCollection(r.Client).FindOne(ctx, bson.M{"_id": objId}).Decode(classroom)
	if err != nil {
		return nil, err
	}
	return classroom, nil
}

// func (r *ClassroomRepositories) GetClassroomsByStudentId(studentId string, page uint, pageLimit uint) (*[]models.Classroom, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(studentId)

// 	findOption := options.Find()
// 	findOption.SetSort(bson.D{
// 		{Key: "type", Value: 1},
// 		{Key: "last_modified", Value: -1},
// 	})
// 	findOption.SetLimit(int64(pageLimit))
// 	findOption.SetSkip(int64((page - 1) * pageLimit))

// 	cursor, err := helpers.ClassroomCollection(r.Client).Find(ctx, bson.M{"_parent_id": objId}, findOption)
// 	if err != nil {
// 		return nil, err
// 	}

// 	directories := &[]models.Classroom{}
// 	cursor.All(ctx, directories)

// 	return directories, nil
// }

func (r *ClassroomRepositories) UpdateClassroom(classroom models.Classroom) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := helpers.ClassroomCollection(r.Client).UpdateByID(ctx, classroom.Id, bson.D{
		{Key: "$set", Value: classroom},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ClassroomRepositories) DeleteClassroom(classroomId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(classroomId)
	_, err := helpers.ClassroomCollection(r.Client).DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return err
	}
	return nil
}
