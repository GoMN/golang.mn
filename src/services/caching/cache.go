// +build !appengine

package caching

import (
	"services"
	"errors"
	"sync"
)

var (
	cache Service
)

func init() {
	cache = Service{}
	cache.data = make(map[string]*interface{})
}

type Service struct {
	lock sync.RWMutex
	data map[string]*interface {}
}

func GetService(rc services.Context) *Service {
	return &cache
}

func (c *Service) Get(key string, item interface {}) (bool, error) {
	var err error
	c.lock.RLock()
	defer c.lock.RUnlock()
	///TODO: not sure this will work
	item, ok := c.data[key]
	if !ok {
		err = errors.New("error retrieving cache key")
	}
	return ok, err
}
func (c *Service) Set(key string, item interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[key] = &item
	return nil
}


