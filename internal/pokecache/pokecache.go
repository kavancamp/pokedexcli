package pokecache

import (
	"time"
	"sync"
)

type cacheEntry struct{
	createdAt time.Time
	val []byte
}

type Cache struct {
	store  map[string]cacheEntry 
	mu sync.Mutex
	interval time.Duration
}
//initialize cache
func NewCache(interval time.Duration) *Cache {
	c := &Cache {
		store: make(map[string]cacheEntry),
		interval: interval,
	}
	go c.reapLoop()
	return c
	
}
//set a value in cache
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
}
//get value from cache
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, found := c.store[key]
	if !found {
		return nil, false
	}
	return entry.val, true
}


func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.store {
			if time.Since(v.createdAt) > c.interval {
				delete(c.store, k)
			}
		}
		c.mu.Unlock()
	}
}
