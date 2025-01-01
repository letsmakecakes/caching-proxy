package cache

import (
	"sync"
	"time"
)

// Entry represents a single cache entry with data, headers, creation time, and content type.
type Entry struct {
	Data        []byte              `json:"data"`
	Headers     map[string][]string `json:"headers"`
	CreatedAt   time.Time           `json:"created_at"`
	ContentType string              `json:"content_type"`
}

// Cache provides a thread-safe map for storing cache entries.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// NewCache initializes and returns a new cache instance.
func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]Entry),
	}
}

// Set adds a new entry to the cache with the given key, data, headers, and content type.
func (c *Cache) Set(key string, data []byte, headers map[string][]string, contentType string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = Entry{
		Data:        data,
		Headers:     headers,
		CreatedAt:   time.Now(),
		ContentType: contentType,
	}
}

// Get retrieves a cache entry by its key. It returns the entry and a boolean indicating if the entry exists.
func (c *Cache) Get(key string) (Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	return entry, exists
}

// Clear removes all entries from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]Entry)
}
