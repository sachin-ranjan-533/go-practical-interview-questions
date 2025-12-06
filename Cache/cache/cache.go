package cache

import (
	"sync"

	"github.com/golang/groupcache/singleflight"
)

type Cache struct {
	data sync.Map
	sf   singleflight.Group
}

func NewCache() *Cache {
	return &Cache{
		data: sync.Map{},
		sf:   singleflight.Group{},
	}
}

func (c *Cache) GetOrCompute(key string, fn func() any) (any, bool) {
	if val, ok := c.data.Load(key); ok {
		return val, true
	}

	val, _ := c.sf.Do(key, func() (any, error) {
		v := fn()
		c.data.Store(key, v)
		return v, nil
	})

	return val, false
}
