package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Directory struct {
	Id           primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	ParentId     primitive.ObjectID `bson:"_parent_id" json:"parent_id"`
	OwnerId      primitive.ObjectID `bson:"_owner_id" json:"owner_id" validate:"required"`
	ClassroomId  primitive.ObjectID `bson:"_classroom_id" json:"classroom_id"`
	Name         string             `bson:"name" json:"name" validate:"required"`
	Type         string             `bson:"type" json:"type" validate:"required,oneof=folder post"`
	Color        *string            `bson:"color" json:"color" validate:"omitempty,oneof=maraschino cayenne maroon plum eggplant grape orchird lavender carnation strawberry bubblegum magenta salmon tangerine cantaloupe banana lemon honeydew lime spring clover fern moss flora sea foam spindrift teal sky turquoise"`
	Description  *string            `bson:"description" json:"description"`
	Files        []File             `bson:"files" json:"files" validate:"required_unless=Type folder,dive,unique=Link"`
	SharedWith   []UserWithAccess   `bson:"shared_with" json:"shared_with" validate:"dive"`
	LastModified primitive.DateTime `bson:"last_modified" json:"last_modified" validate:"required"`
	DateCreated  primitive.DateTime `bson:"date_created" json:"date_created" validate:"required"`
}

type File struct {
	Link string `bson:"link" json:"link" validate:"required,url"`
	Type string `bson:"type" json:"type" validate:"required"`
	Size int64  `bson:"size" json:"size" validate:"required,number"`
}

type UserWithAccess struct {
	UserId primitive.ObjectID `bson:"_user_id" json:"user_id" validate:"required"`
	Access string             `bson:"access" json:"access" validate:"required,oneof=read edit"`
}
