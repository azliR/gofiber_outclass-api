package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `bson:"_id" validate:"required"`
	Name     string             `bson:"name" validate:"required"`
	Email    string             `bson:"email" validate:"required,email"`
	Password string             `bson:"password" validate:"required"`
}

/* MONGO JSON SCHEMA VALIDATION
{
  $jsonSchema: {
    title: 'User Object Validation',
    bsonType: 'object',
    required: [
      'name',
      'email',
      'password'
    ],
    properties: {
      name: {
        bsonType: 'string',
        description: '\'name\' must be a string and is required'
      },
      email: {
        bsonType: 'string',
        description: '\'email\' must be a string and is required'
      },
      password: {
        bsonType: 'string',
        description: '\'password\' must be a string and is required'
      },
    }
  }
}
*/
