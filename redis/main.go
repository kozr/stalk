package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func Init() {
	// Initialize the redis connection
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
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

func userUrlPrefixKey(userId string) string {
	return "user_url:" + userId
}

func UpdateUserUrl(userId string, hashedUrl string) error {
	return saveMessageToRedis(userUrlPrefixKey(userId), hashedUrl, 0)
}
