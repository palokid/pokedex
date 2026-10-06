package pokecache

import "time"

func NewCache(interval time.Duration) *Cache {
	var cache Cache
	cache.entries = make(map[string]cacheEntry, 0)
	go cache.reapLoop(interval)
	return &cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	cache.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
	cache.mu.Unlock()
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entry, ok := cache.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		cache.mu.Lock()

		for key, value := range cache.entries {
			if time.Since(value.createdAt) >= interval {
				delete(cache.entries, key)
			}
		}

		cache.mu.Unlock()
	}
}
