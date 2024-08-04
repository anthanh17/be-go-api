package cache

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

type inMemoryClient struct {
	cache      map[string]any
	cacheMutex *sync.Mutex
	logger     *zap.Logger
}

func NewInMemoryClient(logger *zap.Logger) Cachier {
	return &inMemoryClient{
		cache:      make(map[string]any),
		cacheMutex: new(sync.Mutex),
		logger:     logger,
	}
}

func (m inMemoryClient) Set(_ context.Context, key string, data any, _ time.Duration) error {
	m.cache[key] = data
	return nil
}

func (m inMemoryClient) Get(_ context.Context, key string) (any, error) {
	data, ok := m.cache[key]
	if !ok {
		return nil, ErrCacheMiss
	}

	return data, nil
}

func (m inMemoryClient) AddToSet(_ context.Context, key string, data ...any) error {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	set := m.getSet(key)
	set = append(set, data...)
	m.cache[key] = set
	return nil
}

func (m inMemoryClient) IsDataInSet(_ context.Context, key string, data any) (bool, error) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	set := m.getSet(key)

	for i := range set {
		if set[i] == data {
			return true, nil
		}
	}

	return false, nil
}

func (m inMemoryClient) getSet(key string) []any {
	setValue, ok := m.cache[key]
	if !ok {
		return make([]any, 0)
	}

	// Checks if the value retrieved from cache is a slice of type any.
	set, ok := setValue.([]any)
	if !ok {
		return make([]any, 0)
	}

	return set
}

func (m inMemoryClient) SetNX(ctx context.Context, key string, data any, _ time.Duration) (bool, error) {
	// Check key exits in memory
	_, ok := m.cache[key]

	// Not exits
	if !ok {
		// Set
		m.cache[key] = data
		return true, nil
	}
	// Exits key
	return false, nil
}

func (m inMemoryClient) Del(ctx context.Context, key string) error {
	delete(m.cache, key)
	return nil
}
