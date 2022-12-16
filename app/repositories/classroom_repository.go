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
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type ClassroomRepositories struct {
	*mongo.Client
}

func (r *ClassroomRepositories) CreateClassroom(classroom models.Classroom, classroomMember models.ClassroomMember) error {
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := r.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := helpers.ClassroomCollection(r.Client).InsertOne(sessCtx, classroom)
		if err != nil {
			return nil, err
		}

		_, err = helpers.ClassroomMemberCollection(r.Client).InsertOne(sessCtx, classroomMember)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, txnOptions)

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

func (r *ClassroomRepositories) GetClassroomByClassCode(classCode string) (*models.Classroom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	classroom := &models.Classroom{}

	err := helpers.ClassroomCollection(r.Client).FindOne(ctx, bson.M{"class_code": classCode}).Decode(classroom)
	if err != nil {
		return nil, err
	}
	return classroom, nil
}

func (r *ClassroomRepositories) UpdateClassroom(classroom models.Classroom) error {
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := r.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err = helpers.ClassroomMemberCollection(r.Client).UpdateMany(
			sessCtx,
			bson.M{"_classroom_id": classroom.Id}, bson.M{
				"$set": bson.M{"classroom_name": classroom.Name},
			})
		if err != nil {
			return nil, err
		}

		_, err = helpers.ClassroomCollection(r.Client).UpdateByID(sessCtx, classroom.Id, bson.D{
			{Key: "$set", Value: classroom},
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, txnOptions)

	if err != nil {
		return err
	}
	return nil
}

func (r *ClassroomRepositories) DeleteClassroom(classroomId string) error {
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := r.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	objId, err := primitive.ObjectIDFromHex(classroomId)
	if err != nil {
		return err
	}

	_, err = session.WithTransaction(context.TODO(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := helpers.ClassroomMemberCollection(r.Client).DeleteMany(
			sessCtx, bson.M{"_classroom_id": objId},
		)
		if err != nil {
			return nil, err
		}

		_, err = helpers.DirectoryCollection(r.Client).DeleteMany(
			sessCtx, bson.M{"_classroom_id": objId},
		)
		if err != nil {
			return nil, err
		}

		_, err = helpers.ClassroomCollection(r.Client).DeleteOne(sessCtx, bson.M{"_id": objId})
		if err != nil {
			return nil, err
		}
		return nil, nil
	}, txnOptions)
	if err != nil {
		return err
	}
	return nil
}
