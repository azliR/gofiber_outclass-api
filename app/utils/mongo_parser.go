package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func ToObjectIdString(objectId primitive.ObjectID) *string {
	var objectIdStr *string
	if objectId.IsZero() {
		objectIdStr = nil
	} else {
		objectIdStrStr := objectId.Hex()
		objectIdStr = &objectIdStrStr
	}
	return objectIdStr
}
