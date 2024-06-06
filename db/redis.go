package db

import "github.com/redis/go-redis/v9"

var RedisTypeRepo *redis.Client

func InitRedis() {
	RedisTypeRepo = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
	})
}
func GetRedis() *redis.Client {
	return RedisTypeRepo
}
