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

type UserRepositories struct {
	*mongo.Client
}

func (r *UserRepositories) CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := helpers.UserCollection(r.Client).InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositories) GetUserById(userId string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := &models.User{}

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := helpers.UserCollection(r.Client).FindOne(ctx, bson.M{"_id": objId}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositories) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := &models.User{}
	err := helpers.UserCollection(r.Client).FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
