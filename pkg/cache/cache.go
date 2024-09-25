// pkg/cache/cache.go
package cache

import (
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

type Cache struct {
	cache *lru.Cache[string, Item]
}

type Item struct {
	Value      interface{}
	Expiration int64
}

func New(size int) (*Cache, error) {
	lruCache, err := lru.New[string, Item](size)
	if err != nil {
		return nil, err
	}
	return &Cache{
		cache: lruCache,
	}, nil
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.cache.Add(key, Item{
		Value:      value,
		Expiration: time.Now().Add(duration).UnixNano(),
	})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, found := c.cache.Get(key)
	if !found || time.Now().UnixNano() > item.Expiration {
		return nil, false
	}
	return item.Value, true
}
