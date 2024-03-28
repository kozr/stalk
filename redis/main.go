package redis

import (
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

func SaveMessageToRedis(connectionID, message string) {
}
