package caching

import ()

type Cacher interface {
	Get(key string, item interface {}) (bool, error)
	//GetObj(key string) (interface{}, error)
	Set(key string, item interface {}) error
//	SetObj(key string, item interface{}) error
}

