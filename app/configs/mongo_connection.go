package configs

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"outclass-api/app/repositories"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	mongo.Client
	*repositories.UserRepositories
	*repositories.ClassroomRepositories
	*repositories.ClassroomMemberRepositories
	*repositories.DirectoryRepositories
	*repositories.EventRepositories
	*repositories.TaskRepositories
}

var MongoDb, _ = MongoConnection()

func GetMongoConnection() (*Repositories, error) {
	db := MongoDb
	if db != nil {
		return db, nil
	}
	mongoDb, err := MongoConnection()
	if err != nil {
		return nil, err
	}
	db = mongoDb
	return db, nil
}

func MongoConnection() (*Repositories, error) {
	godotenv.Load()
	mongoOptions := options.Client()
	mongoOptions.ApplyURI(os.Getenv("MONGO_URI"))

	env := os.Getenv("ENV")
	serverScheme := os.Getenv("SERVER_SCHEME")
	if serverScheme == "https" || env == "production" {
		mongoOptions.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	client, err := mongo.NewClient(mongoOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB")

	return &Repositories{
		Client:                      *client,
		UserRepositories:            &repositories.UserRepositories{Client: client},
		ClassroomRepositories:       &repositories.ClassroomRepositories{Client: client},
		ClassroomMemberRepositories: &repositories.ClassroomMemberRepositories{Client: client},
		DirectoryRepositories:       &repositories.DirectoryRepositories{Client: client},
		EventRepositories:           &repositories.EventRepositories{Client: client},
		TaskRepositories:            &repositories.TaskRepositories{Client: client},
	}, nil
}
