package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	Id           primitive.ObjectID  `bson:"_id" validate:"required"`
	CreatorId    primitive.ObjectID  `bson:"_owner_id" validate:"required"`
	ClassroomId  *primitive.ObjectID `bson:"_classroom_id"`
	Title        string              `bson:"title" validate:"required"`
	Details      *string             `bson:"details"`
	Date         primitive.DateTime  `bson:"date" validate:"required"`
	Repeat       string              `bson:"repeat" validate:"required,oneof=none daily weekly monthly yearly"`
	Color        string              `bson:"color" validate:"required,oneof=maraschino cayenne maroon plum eggplant grape orchid lavender carnation strawberry bubblegum magenta salmon tangerine cantaloupe banana lemon honeydew lime spring clover fern moss flora seafoam spindrift teal sky turquoise"`
	Files        []File              `bson:"files" validate:"dive,unique=Link"`
	LastModified primitive.DateTime  `bson:"last_modified" validate:"required"`
	DateCreated  primitive.DateTime  `bson:"date_created" validate:"required"`
}
