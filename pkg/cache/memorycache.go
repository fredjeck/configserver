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

// A basic (but thread safe) memory cache
type MemoryCache struct {
	stop chan struct{}

	waitGroup sync.WaitGroup
	mutex     sync.RWMutex
	entries   map[uint64]cacheEntry
	logger    zap.Logger
}

var (
	// Error thrown when a key cannot be found
	ErrKeyNotInCache = errors.New("the provided key could not be found in the memory cache")
)

// Instantiates a new memory cache which automatically evicts entries after the given retention period
// evicterRunInterval controls the interval at which the cache evicter runs
func NewMemoryCache(evicterRunInterval time.Duration, logger zap.Logger) *MemoryCache {
	cache := &MemoryCache{
		entries: make(map[uint64]cacheEntry),
		stop:    make(chan struct{}),
		logger:  logger,
	}

	cache.waitGroup.Add(1)
	go func(cleanupInterval time.Duration) {
		defer cache.waitGroup.Done()
		cache.startEvicter(cleanupInterval)
	}(evicterRunInterval)

	return cache
}

// Warning : This function must be called witin a goroutine as it will never return until the stop channel is used.
// Starts the cache eviction process
func (cache *MemoryCache) startEvicter(interval time.Duration) {
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

// Forces the eviction of outdated entries
func (cache *MemoryCache) Flush() {
	cache.mutex.Lock()
	evicted := 0
	for id, entry := range cache.entries {
		if entry.expireAtTimestamp <= time.Now().Unix() {
			evicted++
			delete(cache.entries, id)
		}
		cache.logger.Sugar().Infof("%d entries were evicted from cache", evicted)
	}
	cache.mutex.Unlock()
}

// Stops the cache evicter process keeping entries forever in memory
func (cache *MemoryCache) Stop() {
	close(cache.stop)
	cache.waitGroup.Wait()
}

// Updates the value value for the provided key and keeps it in the cache until the provided eviction time is reached
func (cache *MemoryCache) Set(key string, content []byte, eviction time.Time) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.entries[createHash(key)] = cacheEntry{
		content:           content,
		expireAtTimestamp: eviction.Unix(),
	}
}

// Gets the value stored for the provided key or returns ErrKeyNotInCache if the key was not found
func (cache *MemoryCache) Get(key string) ([]byte, error) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	entry, ok := cache.entries[createHash(key)]
	if !ok {
		return nil, ErrKeyNotInCache
	}

	return entry.content, nil
}

// Remove the value stored for the provided key
func (cache *MemoryCache) Remove(key string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	delete(cache.entries, createHash(key))
}

// Creates a has for the provided key
func createHash(path string) uint64 {
	h := fnv.New64()
	h.Write([]byte(path))
	return h.Sum64()
}
