package testutils

import "github.com/go-redis/redis"

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: ServerAddr,
	})
}
