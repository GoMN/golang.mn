package main

import (
	"sync"
)

var (
	cache Cache
)

func init() {
	cache = Cache{}
	cache.data = make(map[string]*[]byte)
	cache.gen = make(map[string]interface {})
}

type Cache struct {
	lock sync.RWMutex
	data map[string]*[]byte
	gen map[string]interface{}
}

/// caching streams
func (c * Cache) Get(key string) *[]byte {
	c.lock.RLock()
	defer c.lock.RUnlock()
	d, _ := c.data[key]
	return d

}
func (c * Cache) Set(key string, d *[]byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[key] = d
}

//caching objects
func (c * Cache) GetGeneric(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	d, _ := c.gen[key]
	return d

}
func (c * Cache) SetGeneric(key string, d interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.gen[key] = d
}
