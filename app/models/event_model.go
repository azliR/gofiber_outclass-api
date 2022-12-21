package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	Id           primitive.ObjectID  `bson:"_id" validate:"required"`
	CreatorId    primitive.ObjectID  `bson:"_owner_id" validate:"required"`
	ClassroomId  *primitive.ObjectID `bson:"_classroom_id"`
	Name         string              `bson:"name" validate:"required"`
	StartDate    primitive.DateTime  `bson:"start_date" validate:"required"`
	EndDate      *primitive.DateTime `bson:"end_date"`
	Repeat       string              `bson:"repeat" validate:"required,oneof=no daily weekly monthly yearly"`
	Color        string              `bson:"color" validate:"required,oneof=maraschino cayenne maroon plum eggplant grape orchid lavender carnation strawberry bubblegum magenta salmon tangerine cantaloupe banana lemon honeydew lime spring clover fern moss flora seafoam spindrift teal sky turquoise"`
	Description  *string             `bson:"description"`
	LastModified primitive.DateTime  `bson:"last_modified" validate:"required"`
	DateCreated  primitive.DateTime  `bson:"date_created" validate:"required"`
}
