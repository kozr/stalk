package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const DEFAULT_TIMEOUT = 5 * time.Second

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

func put(key string, value string) error {
	return putWithExpiration(key, value, 0)
}

func putWithExpiration(key string, value string, expiration time.Duration) error {
	operationCtx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	err := redisClient.Set(operationCtx, key, value, expiration).Err()

	if err != nil {
		log.Printf("Failed to save message to Redis: %v", err)
		return err
	}

	return nil
}

func remove(key string) error {
	operationCtx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	err := redisClient.Del(operationCtx, key).Err()

	if err != nil {
		log.Printf("Failed to remove message from Redis: %v", err)
		return err
	}

	return nil
}

func get(key string) (string, error) {
	operationCtx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	result, err := redisClient.Get(operationCtx, key).Result()

	if err != nil {
		log.Printf("Failed to get message from Redis: %v", err)
		return "", err
	}

	return result, nil
}
