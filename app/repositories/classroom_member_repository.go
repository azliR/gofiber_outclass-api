package repositories

import (
	"context"
	"outclass-api/app/models"
	"outclass-api/app/repositories/helpers"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClassroomMemberRepositories struct {
	*mongo.Client
}

func (r *ClassroomMemberRepositories) CreateClassroomMember(classroomMember models.ClassroomMember) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := helpers.ClassroomMemberCollection(r.Client).InsertOne(ctx, classroomMember)
	if err != nil {
		return err
	}
	return nil
}

func (r *ClassroomMemberRepositories) CountClassroomMembersByClassroomId(classroomId string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(classroomId)

	count, err := helpers.ClassroomMemberCollection(r.Client).CountDocuments(ctx, bson.M{"_classroom_id": objId})
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *ClassroomMemberRepositories) GetClassroomMemberByUserIdAndClassroomId(userId string, classroomId string) (*models.ClassroomMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	classroomMember := &models.ClassroomMember{}

	userIdObj, _ := primitive.ObjectIDFromHex(userId)
	classroomIdObj, _ := primitive.ObjectIDFromHex(classroomId)

	err := helpers.ClassroomMemberCollection(r.Client).FindOne(
		ctx,
		bson.M{
			"_user_id":      userIdObj,
			"_classroom_id": classroomIdObj,
		},
	).Decode(classroomMember)

	if err != nil {
		return nil, err
	}
	return classroomMember, nil
}

func (r *ClassroomMemberRepositories) GetClassroomMembersByClassroomId(classroomId string, page uint, pageLimit uint) (*[]models.ClassroomMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(classroomId)

	findOption := options.Find()
	findOption.SetSort(bson.D{
		{Key: "role", Value: 1},
		{Key: "name", Value: 1},
	})
	findOption.SetLimit(int64(pageLimit))
	findOption.SetSkip(int64((page - 1) * pageLimit))

	cursor, err := helpers.ClassroomMemberCollection(r.Client).Find(ctx, bson.M{"_classroom_id": objId}, findOption)
	if err != nil {
		return nil, err
	}

	classroomMembers := &[]models.ClassroomMember{}
	cursor.All(ctx, classroomMembers)

	return classroomMembers, nil
}

func (r *ClassroomMemberRepositories) GetClassroomMembersByUserId(studentId string, page uint, pageLimit uint) (*[]models.ClassroomMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(studentId)

	findOption := options.Find()
	findOption.SetSort(bson.D{
		{Key: "classroom_name", Value: 1},
	})
	findOption.SetLimit(int64(pageLimit))
	findOption.SetSkip(int64((page - 1) * pageLimit))

	cursor, err := helpers.ClassroomMemberCollection(r.Client).Find(ctx, bson.M{"_user_id": objId}, findOption)
	if err != nil {
		return nil, err
	}

	classroomMembers := &[]models.ClassroomMember{}
	cursor.All(ctx, classroomMembers)

	return classroomMembers, nil
}
