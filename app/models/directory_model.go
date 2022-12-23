package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Directory struct {
	Id           primitive.ObjectID  `bson:"_id" validate:"required"`
	ParentId     *primitive.ObjectID `bson:"_parent_id"`
	OwnerId      primitive.ObjectID  `bson:"_owner_id" validate:"required"`
	ClassroomId  *primitive.ObjectID `bson:"_classroom_id"`
	Name         string              `bson:"name" validate:"required"`
	Type         string              `bson:"type" validate:"required,oneof=folder post"`
	Color        *string             `bson:"color" validate:"omitempty,oneof=maraschino cayenne maroon plum eggplant grape orchid lavender carnation strawberry bubblegum magenta salmon tangerine cantaloupe banana lemon honeydew lime spring clover fern moss flora seafoam spindrift teal sky turquoise"`
	Description  *string             `bson:"description"`
	Files        []File              `bson:"files" validate:"required_unless=Type folder,dive,unique=Link"`
	SharedWith   []UserWithAccess    `bson:"shared_with" validate:"dive"`
	LastModified primitive.DateTime  `bson:"last_modified" validate:"required"`
	DateCreated  primitive.DateTime  `bson:"date_created" validate:"required"`
}

type File struct {
	Id   string  `bson:"id" validate:"required"`
	Name string  `bson:"name" validate:"required"`
	Link string  `bson:"link" validate:"required,url"`
	Type *string `bson:"type"`
	Size *int64  `bson:"size" validate:"omitempty,number"`
}

type UserWithAccess struct {
	UserId primitive.ObjectID `bson:"_user_id" validate:"required"`
	Access string             `bson:"access" validate:"required,oneof=r w rw"`
}

/* MONGO JSON SCHEMA VALIDATION
{
  $jsonSchema: {
    title: 'Directory Object Validation',
    bsonType: 'object',
    required: [
      '_parent_id',
      '_owner_id',
      '_classroom_id',
      'name',
      'type',
      'description',
      'last_modified',
      'date_created'
    ],
    properties: {
      _parent_id: {
        bsonType: 'objectId',
        description: '\'_parent_id\' must be an objectId and is required'
      },
      _owner_id: {
        bsonType: 'objectId',
        description: '\'_owner_id\' must be an objectId and is required'
      },
      _classroom_id: {
        bsonType: 'objectId',
        description: '\'_classroom_id\' must be an objectId and is required'
      },
      name: {
        bsonType: 'string',
        description: '\'name\' must be a string and is required'
      },
      type: {
        bsonType: 'string',
        'enum': [
          'folder',
          'post'
        ],
        description: '\'type\' must be either \'folder\' or \'post\' and is required'
      },
      color: {
        bsonType: [
          'null',
          'string'
        ],
        'enum': [
          null,
          'maraschino',
          'cayenne',
          'maroon',
          'plum',
          'eggplant',
          'grape',
          'orchid',
          'lavender',
          'carnation',
          'strawberry',
          'bubblegum',
          'magenta',
          'salmon',
          'tangerine',
          'cantaloupe',
          'banana',
          'lemon',
          'honeydew',
          'lime',
          'spring',
          'clover',
          'fern',
          'moss',
          'flora',
          'seafoam',
          'spindrift',
          'teal',
          'sky',
          'turquoise'
        ],
        description: '\'color\' must be either \'maraschino\', \'cayenne\', \'maroon\', \'plum\', \'eggplant\', \'grape\', \'orchid\', \'lavender\', \'carnation\', \'strawberry\', \'bubblegum\', \'magenta\', \'salmon\', \'tangerine\', \'cantaloupe\', \'banana\', \'lemon\', \'honeydew\', \'lime\', \'spring\', \'clover\', \'fern\', \'moss\', \'flora\', \'seafoam\', \'spindrift\', \'teal\', \'sky\', or \'turquoise\''
      },
      description: {
        bsonType: [
          'null',
          'string'
        ],
        description: '\'description\' must be a string and is required'
      },
      files: {
        bsonType: [
          'null',
          'array'
        ],
        uniqueItems: true,
        description: '\'files\' must be an array',
        items: {
          bsonType: 'object',
          required: [
            'link'
          ],
          properties: {
            link: {
              bsonType: 'string',
              description: '\'link\' must be a string and is required'
            },
            type: {
              bsonType: [
                'null',
                'string'
              ],
              description: '\'type\' must be a string and is required'
            },
            size: {
              bsonType: [
                'null',
                'number'
              ],
              description: '\'size\' must be a number and is required'
            }
          }
        }
      },
      shared_with: {
        bsonType: [
          'null',
          'array'
        ],
        description: '\'shared_with\' must be an array',
        items: {
          bsonType: 'object',
          required: [
            '_user_id',
            'access'
          ],
          properties: {
            _user_id: {
              bsonType: 'objectId',
              description: '\'_user_id\' must be an objectId and is required'
            },
            access: {
              bsonType: 'string',
              'enum': [
                'read',
                'edit'
              ],
              description: '\'access\' must be either \'read\' or \'edit\' and is required'
            }
          }
        }
      },
      last_modified: {
        bsonType: 'date',
        description: '\'last_modified\' must be a date and is required'
      },
      date_created: {
        bsonType: 'date',
        description: '\'date_created\' must be a date and is required'
      }
    }
  }
}
*/
