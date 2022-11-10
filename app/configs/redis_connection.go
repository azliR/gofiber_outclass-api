package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
)

var RedisDb = RedisConnection()

func RedisConnection() *redis.Client {
	godotenv.Load()
	dbNumber, err := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))

	if err != nil {
		log.Fatal(err)
	}

	redisConnURL := fmt.Sprintf(
		"%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)

	options := &redis.Options{
		Addr:     redisConnURL,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	}
	fmt.Println("Connected to Redis")
	return redis.NewClient(options)
}
