package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Classroom struct {
	Id          primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	ClassCode   string             `bson:"class_code" json:"class_code" validate:"required"`
	Description *string            `bson:"description" json:"description"`
}
