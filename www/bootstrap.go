package main

import (
	"appengine"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	boot      = bootstrapper{}
	meetupSvc = meetupService{}
)

const BOOTSTRAP_KEY = "bootstrap"

type bootstrapper struct{
	context     appengine.Context
	initialized bool
	Bootstrap   bootstrap `json:"bootstrap"`
}

///models
type bootstrap struct{
	Members       []Member `json:"members"`
	MemberCoords  []memberCoord `json:"member_coords"`
}
type memberCoord struct {
	Title string `json:"title"`
	Lat   float32 `json:"lat"`
	Lon   float32 `json:"lon"`
}

func (b *bootstrapper) Scope(r *http.Request) {
	b.context = appengine.NewContext(r)
	if !b.initialized {
		b.initialize()
	}
}

func (b *bootstrapper) initialize() error {
	if b.initialized {
		log.Println("cached bootstrap will be used")
		return nil
	}
	meetupSvc.SetContext(b.context)
	var wg sync.WaitGroup
	cbg := cache.GetGeneric(BOOTSTRAP_KEY)
	test, ok := cbg.(bootstrap)

	if ok {
		log.Println("cached bootstrap reset")
		b.Bootstrap = test
	}else {
		log.Println("bootstrap reinitializing")
		wg.Add(1)
		go func(boot *bootstrap, svc meetupService) {
			defer wg.Done()
			members, _ := svc.getMembers()
			boot.Members = members

			for _, m := range members {
				boot.MemberCoords = append(boot.MemberCoords, memberCoord{
					m.Name,
					m.Lat,
					m.Lon,
				})
			}

		}(&b.Bootstrap, meetupSvc)

		wg.Add(1)
		go func(boot *bootstrap) {
			defer wg.Done()

		}(&b.Bootstrap)

		//wait for everything to bootstrap or fail
		wg.Wait()

		//cache this result
		cache.SetGeneric(BOOTSTRAP_KEY, b.Bootstrap)
		b.initialized = true
		appdata.Bootstrap = boot.Bootstrap;
		//fire and forget the cache timeout
		go func(timeout int64) {
			time.Sleep(time.Duration(timeout) * time.Millisecond)
			clearBootstrapCache()
			b.refresh()
		}(config.Cache.LocalTimeout)
	}
	return nil
}

func (b *bootstrapper) refresh() error {
	//zap cache and initialize
	clearBootstrapCache()
	b.initialized = false
	return b.initialize()
}

/// zap the bootstrap
func clearBootstrapCache() {
	cache.SetGeneric(BOOTSTRAP_KEY, nil)
}
