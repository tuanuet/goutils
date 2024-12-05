//go:build wireinject
// +build wireinject

package redis

import (
	"context"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

func InitializeRedis(ctx context.Context, dns RedisConnectionString) *redis.Client {
	wire.Build(newRedisClient)
	return nil
}
