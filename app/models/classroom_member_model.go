package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassroomMember struct {
	Id            primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	UserId        primitive.ObjectID `bson:"_user_id" json:"user_id" validate:"required"`
	ClassroomId   primitive.ObjectID `bson:"_classroom_id" json:"classroom_id" validate:"required"`
	StudentId     string             `bson:"student_id" json:"student_id" validate:"required"`
	ClassroomName string             `bson:"classroom_name" json:"classroom_name" validate:"required"`
	Name          string             `bson:"name" json:"name" validate:"required"`
	Role          uint16             `bson:"role" json:"role" validate:"required,number,oneof=1 2 3"`
}

// func ToClassroomMember(classroomMemberDto dtos.CreateClassroomDto) ClassroomMember {
// }
