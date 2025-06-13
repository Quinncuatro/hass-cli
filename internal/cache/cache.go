package cache

import (
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
	Size() int
}

type MemoryCache struct {
	mu     sync.RWMutex
	items  map[string]*cacheItem
	maxTTL time.Duration
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

func NewMemoryCache(maxTTL time.Duration) *MemoryCache {
	cache := &MemoryCache{
		items:  make(map[string]*cacheItem),
		maxTTL: maxTTL,
	}
	
	go cache.cleanup()
	
	return cache
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}
	
	if time.Now().After(item.expiresAt) {
		delete(c.items, key)
		return nil, false
	}
	
	return item.value, true
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if ttl > c.maxTTL {
		ttl = c.maxTTL
	}
	
	c.items[key] = &cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.items, key)
}

func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items = make(map[string]*cacheItem)
}

func (c *MemoryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return len(c.items)
}

func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}