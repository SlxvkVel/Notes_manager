package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func InitRedis() {
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "", 
		DB:       0,  
	})


	var err error
	for i := 0; i < 5; i++ {
		_, err = redisClient.Ping(ctx).Result()
		if err == nil {
			break
		}
		log.Printf("Redis connection attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Printf("Failed to connect to Redis after 5 attempts: %v", err)
		return
	}

	log.Println("✅ Connected to Redis successfully")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func StoreUserSession(userID int32, token string, expiration time.Duration) error {
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("session:%d", userID)
	err := redisClient.Set(ctx, key, token, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to store session in Redis: %w", err)
	}
	return nil
}

func GetUserSession(userID int32) (string, error) {
	if redisClient == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("session:%d", userID)
	token, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get session from Redis: %w", err)
	}
	return token, nil
}


func DeleteUserSession(userID int32) error {
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}

	key := fmt.Sprintf("session:%d", userID)
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session from Redis: %w", err)
	}
	return nil
}