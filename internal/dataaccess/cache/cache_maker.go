package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anthanh17/be-go-api/configs"

	"go.uber.org/zap"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

const (
	CacheTypeInMemory configs.CacheType = "in_memory"
	CacheTypeRedis    configs.CacheType = "redis"
)

// Repository pattern
type Cachier interface {
	Set(ctx context.Context, key string, data any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)

	// Adds one or more values ​​to a set
	AddToSet(ctx context.Context, key string, data ...any) error
	IsDataInSet(ctx context.Context, key string, data any) (bool, error)

	// Using `SETNX`: SET if Not Exists. A.K.A `SET if Not Exists.`
	SetNX(ctx context.Context, key string, data any, ttl time.Duration) (bool, error)

	Del(ctx context.Context, key string) error
}

// Factory pattern
func NewCachierClient(
	cacheConfig configs.CacheConfig,
	logger *zap.Logger,
) (Cachier, error) {
	switch cacheConfig.Type {
	case CacheTypeInMemory:
		return NewInMemoryClient(logger), nil

	case CacheTypeRedis:
		return NewRedisClient(cacheConfig, logger), nil

	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheConfig.Type)
	}
}
