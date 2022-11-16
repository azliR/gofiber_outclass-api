package helpers

import "go.mongodb.org/mongo-driver/mongo"

func getCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("outclass").Collection(collectionName)
	return collection
}

func UserCollection(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "users")
}

func ClassroomCollection(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "classrooms")
}

func ClassroomMemberCollection(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "classroom_members")
}

func DirectoryCollection(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "directories")
}
