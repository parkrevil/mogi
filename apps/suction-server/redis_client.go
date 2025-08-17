package main

import (
	"context"
	"time"

	"go/common"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RedisClient struct {
	client *redis.Client
	logger *zap.Logger
}

func NewRedisClient(logger *zap.Logger, config *common.Config, lifecycle fx.Lifecycle) (*RedisClient, error) {
	opts, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	redisClient := &RedisClient{
		client: client,
		logger: logger,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to Redis",
				zap.String("address", opts.Addr),
				zap.Int("db", opts.DB),
				zap.Int("pool_size", opts.PoolSize),
			)

			ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()

			if err := client.Ping(ctx).Err(); err != nil {
				return err
			}

			logger.Info("Successfully connected to Redis")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Redis connection")

			return client.Close()
		},
	})

	return redisClient, nil
}
