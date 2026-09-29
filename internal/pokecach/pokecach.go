package pokecach

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}
type Cache struct {
	cache map[string]cacheEntry
	mu    sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cache: make(map[string]cacheEntry),
	}
	go cache.reaploop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}
func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	} else {
		return entry.val, true
	}
}
func (c *Cache) reap(t time.Duration) {
	for k, entry := range c.cache {
		if time.Since(entry.createdAt) > t {
			delete(c.cache, k)
		}
	}
}

func (c *Cache) reaploop(t time.Duration) {
	for true {
		time.Sleep(t)
		c.reap(t)
	}
}
