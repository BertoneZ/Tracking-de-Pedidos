package db

import (
	//"context"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return rdb
}