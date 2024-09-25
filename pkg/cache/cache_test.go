package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cacheSize := 2
	cache, err := New(cacheSize)
	assert.NoError(t, err)

	// Test setting and getting an item
	cache.Set("key1", "value1", 5*time.Second)
	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "value1", value)

	// Test expiration of an item
	cache.Set("key2", "value2", 1*time.Second)
	time.Sleep(2 * time.Second)
	value, found = cache.Get("key2")
	assert.False(t, found)
	assert.Nil(t, value)

	// Test cache size limit
	cache.Set("key3", "value3", 5*time.Second)
	cache.Set("key4", "value4", 5*time.Second)
	_, found = cache.Get("key1")
	assert.False(t, found) // key1 should be evicted due to cache size limit
}
