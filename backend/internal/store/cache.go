package store

import (
	"sync"
	"time"
)

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

type CacheStore struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func NewCacheStore() *CacheStore {
	cache := &CacheStore{
		items: make(map[string]cacheItem),
	}
	go cache.startCleanup()
	return cache
}

func (c *CacheStore) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(duration),
	}
}

func (c *CacheStore) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}

func (c *CacheStore) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *CacheStore) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		c.cleanup()
	}
}

func (c *CacheStore) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.expiration) {
			delete(c.items, key)
		}
	}
}
