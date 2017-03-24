package redis

import (
	"config"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func init() {
	redisHost := config.Get("REDIS_HOST", "127.0.0.1")
	redisPort := config.Get("REDIS_PORT", "6379")
	redisPass := config.Get("REDIS_PASSWORD", "")
	addr := redisHost + ":" + redisPort
	redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPass,
		DB:       0,
	})
}
