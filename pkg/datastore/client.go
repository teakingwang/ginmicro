package datastore

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/teakingwang/ginmicro/config"
)

type redisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient 创建一个 Redis 封装实例
func NewRedisClient() Store {
	cfg := config.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(fmt.Sprintf("Redis connect error: %v", err))
	}

	return &redisClient{
		client: rdb,
		ctx:    context.Background(),
	}
}
