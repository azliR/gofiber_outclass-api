package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"outclass-api/app/repositories"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repositories struct {
	*repositories.UserRepositories
	*repositories.DirectoryRepositories
}

var MongoDb = MongoConnection()

func MongoConnection() *Repositories {
	godotenv.Load()
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil
	}

	fmt.Println("Connected to MongoDB")

	return &Repositories{
		UserRepositories:      &repositories.UserRepositories{Client: client},
		DirectoryRepositories: &repositories.DirectoryRepositories{Client: client},
	}
}
