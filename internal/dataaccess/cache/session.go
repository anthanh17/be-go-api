package cache

import (
	"context"
	"fmt"

	"github.com/anthanh17/be-go-api/internal/utils"
	"go.uber.org/zap"
)

type SessionCache interface {
	GetSession(ctx context.Context, id string) (SessionType, error)
	SetSession(ctx context.Context, id string, data SessionType) error

	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, data any) error
	// SetNX
	SetPingLock(ctx context.Context, key string, data any) (bool, error)
	Del(ctx context.Context, key string) error
}

type SessionType struct {
	SessionID string `json:"sessionId"`
	Username  string `json:"username"`
}

type sessionCache struct {
	cachier Cachier
	logger  *zap.Logger
}

func NewSessionCache(cachier Cachier, logger *zap.Logger) SessionCache {
	return &sessionCache{
		cachier: cachier,
		logger:  logger,
	}
}

func (s sessionCache) getSessionCacheKey(id string) string {
	return fmt.Sprintf("session_cache_key:%s", id)
}

func (s sessionCache) GetSession(ctx context.Context, id string) (SessionType, error) {
	logger := utils.LoggerWithContext(ctx, s.logger).With(zap.String("id", id))

	// get key cache
	cacheKey := s.getSessionCacheKey(id)

	// Get cache data
	cacheEntry, err := s.cachier.Get(ctx, cacheKey)
	if err != nil {
		// logger.With(zap.Error(err)).Error("failed to get session key cache")
		logger.Info("failed to get session key cache")
		return SessionType{}, err
	}

	// If miss cache
	if cacheEntry == nil {
		return SessionType{}, ErrCacheMiss
	}

	// Check data type session
	sessionData, ok := cacheEntry.(SessionType)
	if !ok {
		logger.Error("cache entry is not of SessionType")
		return SessionType{}, nil
	}

	return sessionData, nil
}

func (s sessionCache) Get(ctx context.Context, key string) (string, error) {
	logger := utils.LoggerWithContext(ctx, s.logger).With(zap.String("key", key))

	// Get cache data
	cacheEntry, err := s.cachier.Get(ctx, key)
	if err != nil {
		logger.Info("failed to get session key cache")
		return "", err
	}

	// If miss cache
	if cacheEntry == nil {
		logger.Info("miss cache")
		return "", ErrCacheMiss
	}

	// Check data type session
	cacheData, ok := cacheEntry.(string)
	if !ok {
		logger.Error("cache entry is not of string")
		return "", nil
	}

	return cacheData, nil
}

func (s sessionCache) SetSession(ctx context.Context, id string, data SessionType) error {
	logger := utils.LoggerWithContext(ctx, s.logger).With(zap.String("id", id))

	// get key cache
	cacheKey := s.getSessionCacheKey(id)

	if err := s.cachier.Set(ctx, cacheKey, data, 0); err != nil {
		// logger.With(zap.Error(err)).Error("failed to insert token public key into cache")
		logger.Info("failed to insert token public key into cache")
		return err
	}

	return nil
}

func (s sessionCache) Set(ctx context.Context, key string, data any) error {
	logger := utils.LoggerWithContext(ctx, s.logger).With(zap.String("key", key))

	if err := s.cachier.Set(ctx, key, data, 0); err != nil {
		logger.Info("failed to insert token public key into cache")
		return err
	}

	return nil
}

func (s sessionCache) SetPingLock(ctx context.Context, key string, data any) (bool, error) {
	ok, err := s.cachier.SetNX(ctx, key, data, 0)
	if err != nil {
		s.logger.Info("failed SetPingLock")
		return false, err
	}

	return ok, nil
}

func (s sessionCache) Del(ctx context.Context, key string) error {
	err := s.cachier.Del(ctx, key)
	if err != nil {
		s.logger.Info("failed to delete session key cache")
		return err
	}
	return nil
}
