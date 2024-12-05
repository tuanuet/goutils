package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
)

var (
	db *redis.Client
)

type RedisConnectionString string

func newRedisClient(ctx context.Context, dns RedisConnectionString) *redis.Client {
	var once sync.Once
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr: string(dns),
		})

		db = rdb
	})
	return db
}
