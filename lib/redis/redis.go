package redis

import (
    "os"
    "sync"

    "github.com/redis/go-redis/v9"
)

var (
    clientInstance *redis.Client
    once           sync.Once
)

// GetClient returns a singleton Redis client instance.
func GetClient() *redis.Client {
    once.Do(func() {
        redisURL := os.Getenv("REDIS_URL")
        if redisURL == "" {
            // Default URL if REDIS_URL is not set (e.g., no auth)
            redisURL = "redis://localhost:6379/1"
        }
        options, err := redis.ParseURL(redisURL)
        if err != nil {
            panic(err)
        }
        clientInstance = redis.NewClient(options)
    })
    return clientInstance
}