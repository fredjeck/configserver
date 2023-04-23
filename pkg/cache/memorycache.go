package cache

import (
	"errors"
	"hash/fnv"
	"sync"
	"time"

	"go.uber.org/zap"
)

type cacheEntry struct {
	content           []byte
	expireAtTimestamp int64
}

// MemoryCache is a basic (but thread safe) memory cache
type MemoryCache struct {
	stop      chan struct{}
	waitGroup sync.WaitGroup
	mutex     sync.RWMutex
	entries   map[uint64]cacheEntry
}

var (
	ErrKeyNotInCache = errors.New("the provided key could not be found in the memory cache")
)

// NewMemoryCache Instantiates a new memory cache which automatically evicts entries after the given retention period
// evictorRunInterval controls the interval at which the cache evictor runs
func NewMemoryCache(evictorRunInterval time.Duration, logger *zap.Logger) *MemoryCache {
	cache := &MemoryCache{
		entries: make(map[uint64]cacheEntry),
		stop:    make(chan struct{}),
	}

	cache.waitGroup.Add(1)
	go func(interval time.Duration) {
		defer cache.waitGroup.Done()
		logger.Sugar().Infof("Memory cache created (cache evictor interval set to %d seconds)", interval)
		cache.startEvictor(interval)
	}(evictorRunInterval)

	return cache
}

// Warning : This function must be called within a goroutine as it will never return until the stop channel is used.
// Starts the cache eviction process
func (cache *MemoryCache) startEvictor(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-cache.stop:
			return
		case <-t.C:
			cache.Flush()
		}
	}
}

// Flush forces the eviction of outdated entries
func (cache *MemoryCache) Flush() {
	cache.mutex.Lock()
	evicted := 0
	for id, entry := range cache.entries {
		if entry.expireAtTimestamp <= time.Now().Unix() {
			evicted++
			delete(cache.entries, id)
		}
	}
	cache.mutex.Unlock()
}

// Stop halts the cache evictor process keeping entries forever in memory
func (cache *MemoryCache) Stop() {
	close(cache.stop)
	cache.waitGroup.Wait()
}

// Set stores the value for the provided key in the cache until the expiry time is reached
func (cache *MemoryCache) Set(key string, content []byte, expiration time.Time) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.entries[createHash(key)] = cacheEntry{
		content:           content,
		expireAtTimestamp: expiration.Unix(),
	}
}

// Get retrieves the value stored for the provided key or returns ErrKeyNotInCache if the key was not found
func (cache *MemoryCache) Get(key string) ([]byte, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	entry, ok := cache.entries[createHash(key)]
	if !ok {
		return nil, ErrKeyNotInCache
	}

	return entry.content, nil
}

// Remove deletes the value stored for the provided key
func (cache *MemoryCache) Remove(key string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	delete(cache.entries, createHash(key))
}

// Creates a hash for the provided key
func createHash(key string) uint64 {
	h := fnv.New64()
	_, _ = h.Write([]byte(key))
	return h.Sum64()
}
