package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Classroom struct {
	Id          primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	ClassCode   string             `bson:"class_code" json:"class_code" validate:"required"`
	Description *string            `bson:"description" json:"description"`
}

/* MONGO JSON SCHEMA VALIDATION
{
  $jsonSchema: {
    title: 'Classroom Object Validation',
    bsonType: 'object',
    required: [
      'name',
      'class_code',
      'description'
    ],
    properties: {
      name: {
        bsonType: 'string',
        description: '\'name\' must be a string and is required'
      },
      class_code: {
        bsonType: 'string',
        description: '\'class_code\' must be a string and is required'
      },
      description: {
        bsonType: [
          'null',
          'string'
        ],
        description: '\'description\' must be a string or null and is required'
      }
    }
  }
}
*/
