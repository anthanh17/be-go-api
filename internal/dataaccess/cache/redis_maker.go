package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anthanh17/be-go-api/configs"
	"github.com/anthanh17/be-go-api/internal/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type redisClient struct {
	redisClient *redis.Client
	logger      *zap.Logger
}

func NewRedisClient(
	cacheConfig configs.CacheConfig,
	logger *zap.Logger,
) Cachier {
	// 127.0.0.1:6379
	addrString := fmt.Sprintf("%s:%d", cacheConfig.Host, cacheConfig.Port)

	return &redisClient{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     addrString,
			Username: cacheConfig.Username,
			Password: cacheConfig.Password,
		}),
		logger: logger,
	}
}

func (r redisClient) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	logger := utils.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data)).
		With(zap.Duration("ttl", ttl))

	if err := r.redisClient.Set(ctx, key, data, ttl).Err(); err != nil {
		logger.With(zap.Error(err)).Error("failed to set data into cache")
		return status.Error(codes.Internal, "failed to set data into cache")
	}

	return nil
}

func (r redisClient) Get(ctx context.Context, key string) (any, error) {
	logger := utils.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key))

	data, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}

		logger.With(zap.Error(err)).Error("failed to get data from cache")
		return nil, status.Error(codes.Internal, "failed to get data from cache")
	}

	return data, nil
}

func (r redisClient) AddToSet(ctx context.Context, key string, data ...any) error {
	logger := utils.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data))

	if err := r.redisClient.SAdd(ctx, key, data...).Err(); err != nil {
		logger.With(zap.Error(err)).Error("failed to set data into set inside cache")
		return status.Error(codes.Internal, "failed to set data into set inside cache")
	}

	return nil
}

func (r redisClient) IsDataInSet(ctx context.Context, key string, data any) (bool, error) {
	logger := utils.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data))

	result, err := r.redisClient.SIsMember(ctx, key, data).Result()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to check if data is member of set inside cache")
		return false, status.Error(codes.Internal, "failed to check if data is member of set inside cache")
	}

	return result, nil
}
