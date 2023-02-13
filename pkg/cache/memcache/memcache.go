package memcache

import (
	"sync"
	"time"
	"where-do-i-sit/pkg/cache"
)

type MemoryCache struct {
	cache map[string]cacheValue
	mu    sync.RWMutex
}

type cacheValue struct {
	value     any
	timestamp time.Time
	expire    time.Duration
}

func (m *MemoryCache) Get(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(value.timestamp.Add(value.expire)) {
		delete(m.cache, key)
		return nil, false
	}

	return value.value, exists
}

func (m *MemoryCache) Set(key string, value any, expire time.Duration) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.cache[key] = cacheValue{
		value:     value,
		timestamp: time.Now(),
		expire:    expire,
	}
}

func NewMemCache() cache.Cache {
	return &MemoryCache{
		cache: make(map[string]cacheValue),
		mu:    sync.RWMutex{},
	}
}

func (m *MemoryCache) String() string {
	return ""
}
