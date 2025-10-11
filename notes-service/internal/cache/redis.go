package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"notes-service/internal/models"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("❌ Failed to connect to Redis: %v", err)
		return
	}
	log.Println("✅ Notes service connected to Redis")
}


func CacheUserNotes(userID int32, notes []models.Note, expiration time.Duration) error {
	if redisClient == nil {
		return nil 
	}

	key := fmt.Sprintf("user:%d:notes", userID)
	jsonData, err := json.Marshal(notes)
	if err != nil {
		return err
	}

	err = redisClient.Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		log.Printf("Warning: failed to cache notes: %v", err)
	}
	return err
}


func GetCachedUserNotes(userID int32) ([]models.Note, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("redis not available")
	}

	key := fmt.Sprintf("user:%d:notes", userID)
	jsonData, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var notes []models.Note
	err = json.Unmarshal([]byte(jsonData), &notes)
	return notes, err
}

func InvalidateUserCache(userID int32) error {
	if redisClient == nil {
		return nil
	}

	key := fmt.Sprintf("user:%d:notes", userID)
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Warning: failed to invalidate cache: %v", err)
	}
	return err
}