package helpers

import "go.mongodb.org/mongo-driver/mongo"

func getCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("outclass").Collection(collectionName)
	return collection
}

func UserCollection(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "users")
}

func DirectoryCollection(client *mongo.Client) *mongo.Collection {
	return getCollection(client, "directories")
}
