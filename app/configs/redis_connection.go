package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
)

var RedisDb, _ = RedisConnection()

func GetRedisConnection() (*redis.Client, error) {
	db := RedisDb
	if db == nil {
		RedisDb, err := RedisConnection()
		if err != nil {
			return nil, err
		}
		fmt.Println("Connected to Redis")
		db = RedisDb
	}
	return db, nil
}

func RedisConnection() (*redis.Client, error) {
	godotenv.Load()
	dbNumber, err := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	if err != nil {
		return nil, err
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisConnURL := redisHost + ":" + redisPort

	options := &redis.Options{
		Addr:     redisConnURL,
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	}
	return redis.NewClient(options), nil
}
