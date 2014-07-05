// +build appengine

package caching

import (
	"appengine"
	"appengine/memcache"
	"services"
)

var (
	cache = Service{}
)

func GetService(rc services.Context) *Service {
	cache.context = appengine.NewContext(rc.Request)
	return &cache
}

type Service struct {
	context appengine.Context
}

func (c *Service) Get(key string, item interface{}) (bool, error) {
	// Get the item from the memcache
	if _, err := memcache.Gob.Get(c.context, key, item); err == memcache.ErrCacheMiss {
		c.context.Infof("item not in the cache")
		return false, err
	} else if err != nil {
		c.context.Errorf("error getting item: %v", err)
		return false, err
	}else {
		c.context.Infof("item found in cache")
		return true, nil
	}
}
func (c *Service) Set(key string, item interface{}) error {
	// Create an Item
	citem := &memcache.Item{
		Key:   key,
		Object: item,
	}
	// Set the item, unconditionally
	if err := memcache.Gob.Set(c.context, citem); err != nil {
		c.context.Errorf("error setting item: %v", err)
		return err
	}
	return nil
}
