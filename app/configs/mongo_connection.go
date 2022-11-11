package configs

import (
	"context"
	"fmt"
	"os"
	"outclass-api/app/repositories"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	*repositories.UserRepositories
	*repositories.ClassroomRepositories
	*repositories.DirectoryRepositories
}

var MongoDb, _ = MongoConnection()

func GetMongoConnection() (*Repositories, error) {
	db := MongoDb
	if db == nil {
		mongoDb, err := MongoConnection()
		if err != nil {
			return nil, err
		}
		db = mongoDb
	}
	return db, nil
}

func MongoConnection() (*Repositories, error) {
	godotenv.Load()
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
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
		UserRepositories:      &repositories.UserRepositories{Client: client},
		ClassroomRepositories: &repositories.ClassroomRepositories{Client: client},
		DirectoryRepositories: &repositories.DirectoryRepositories{Client: client},
	}, nil
}
