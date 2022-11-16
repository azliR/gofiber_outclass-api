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

type UserRepositories struct {
	*mongo.Client
}

func (r *UserRepositories) CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := helpers.UserCollection(r.Client).InsertOne(ctx, user); err != nil {
		return err
	}
	return nil
}

func (r *UserRepositories) GetUserById(userId string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := &models.User{}

	objId, _ := primitive.ObjectIDFromHex(userId)

	if err := helpers.UserCollection(r.Client).FindOne(ctx, bson.M{"_id": objId}).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositories) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := &models.User{}
	if err := helpers.UserCollection(r.Client).FindOne(ctx, bson.M{"email": email}).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositories) UpdateUserById(userId string, updateName bool, user *models.User) error {
	wc := writeconcern.New(writeconcern.WMajority())
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := r.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	objId, _ := primitive.ObjectIDFromHex(userId)

	if _, err := session.WithTransaction(context.TODO(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		if updateName {
			if _, err := helpers.ClassroomMemberCollection(r.Client).UpdateMany(
				sessCtx,
				bson.M{"_user_id": objId},
				bson.M{"$set": bson.M{"name": user.Name}},
			); err != nil {
				return nil, err
			}
		}

		if _, err := helpers.UserCollection(r.Client).UpdateOne(
			sessCtx,
			bson.M{"_id": objId},
			bson.M{"$set": user},
		); err != nil {
			return nil, err
		}

		return nil, nil
	}, txnOptions); err != nil {
		return err
	}

	return nil
}
