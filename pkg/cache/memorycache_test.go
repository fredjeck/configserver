package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestStoreAndRetrieve(t *testing.T) {

	logger, _ := zap.NewDevelopment()
	cache := NewMemoryCache(60*time.Minute, logger)

	key := "random key"
	value := "go is a fantastic language"
	cache.Set(key, []byte(value), time.Now().Add(60*time.Minute))

	result, err := cache.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, value, string(result))
}

func TestEviction(t *testing.T) {

	logger, _ := zap.NewDevelopment()
	cache := NewMemoryCache(1*time.Second, logger)

	key := "random key"
	value := "go is a fantastic language"
	cache.Set(key, []byte(value), time.Now().Add(1*time.Second))
	time.Sleep(1 * time.Second)
	cache.Flush()
	_, err := cache.Get(key)
	assert.ErrorIs(t, err, ErrKeyNotInCache)
}
