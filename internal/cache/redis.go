package cache

import (
	"sql_sharding_engine/internal/config"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	config.Redis = client
}
