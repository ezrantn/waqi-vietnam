package main

import (
	"strings"
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

const (
	// Cache durations
	DefaultCacheDuration = 5 * time.Minute
	LongCacheDuration    = 30 * time.Minute

	// Cache keys
	AirQualityKeyPrefix = "air_quality:"
)

func NewInMemoryCache() *Cache {
	c := &Cache{
		items: make(map[string]CacheItem),
	}

	// Start the cleanup routine
	go c.startCleanupTimer()

	return c
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(duration).UnixNano()
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if item has expired
	if time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Cleanup expired items every minute
func (c *Cache) startCleanupTimer() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		for key, item := range c.items {
			if now > item.Expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// GenerateAirQualityCacheKey creates a standardized cache key for air quality data
func GenerateAirQualityCacheKey(city string) string {
	return AirQualityKeyPrefix + strings.ToLower(city)
}
