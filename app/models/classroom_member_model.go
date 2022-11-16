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

/* MONGO JSON SCHEMA VALIDATION
{
  $jsonSchema: {
    title: 'ClassroomMember Object Validation',
    bsonType: 'object',
    required: [
      '_user_id',
      '_classroom_id',
      'student_id',
      'classroom_name',
      'name',
      'role'
    ],
    properties: {
      _user_id: {
        bsonType: 'objectId',
        description: '\'_user_id\' must be an objectId and is required'
      },
      _classroom_id: {
        bsonType: 'objectId',
        description: '\'_classroom_id\' must be an objectId and is required'
      },
      student_id: {
        bsonType: 'string',
        description: '\'student_id\' must be a string and is required'
      },
      classroom_name: {
        bsonType: 'string',
        description: '\'classroom_name\' must be a string and is required'
      },
      name: {
        bsonType: 'string',
        description: '\'name\' must be a string and is required'
      },
      role: {
        bsonType: 'int',
        'enum': [
          1,
          2,
          3
        ],
        description: '\'role\' must be one of [1, 2, 3] and is required'
      }
    }
  }
}
*/
