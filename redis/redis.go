package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func Init() error {
	// Initialize the redis connection
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("failed to connect to Redis: %v", err)
		return err
	}

	return nil
}

func saveMessageToRedis(key string, value string, expiration time.Duration) error {
	const operationTimeout = 5 * time.Second

	// Create a new context with the operation timeout
	operationCtx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	err := redisClient.Set(operationCtx, key, value, expiration).Err()

	if err != nil {
		log.Printf("Failed to save message to Redis: %v", err)
		return err // Propagate the error for further handling
	}

	return nil
}
