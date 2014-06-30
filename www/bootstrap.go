package main

import (
	"appengine"
	"log"
	"net/http"
	"strconv"
	"sync"
	"text/template"
	"time"
)

var (
	boot        = bootstrapper{}
	meetupSvc   = meetupService{}
	bootFuncMap = template.FuncMap{
	"marshall": marshall, }
	bootTmpl    = template.Must(template.New("bootstrap").Funcs(bootFuncMap).Parse("window.$app.bootstrap = {{ marshall . }};"))

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
	Calendar      Calendar `json:"calendar"`
	Version       string `json:"version"`
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
		b.Bootstrap.Version = appdata.Version
		wg.Add(1)
		go func(boot *bootstrap, svc meetupService) {
			defer wg.Done()
			members, _ := svc.getMembers()

			for _, m := range members {
				boot.Members = append(boot.Members, Member{
					m.ID, m.Joined, m.Bio, m.Link, m.Name, m.City, m.State, m.Photo, m.Other,
				})
				boot.MemberCoords = append(boot.MemberCoords, memberCoord{
					"gopher",
					m.Lat,
					m.Lon,
				})
			}

			boot.Calendar = svc.getMembersCalendar(boot.Members)

		}(&b.Bootstrap, meetupSvc)

//		wg.Add(1)
//		go func(boot *bootstrap) {
//			defer wg.Done()
//			//boot.Calendar =
//
//		}(&b.Bootstrap)

		//wait for everything to bootstrap or fail
		wg.Wait()

		//cache this result
		cache.SetGeneric(BOOTSTRAP_KEY, b.Bootstrap)
		b.initialized = true

		//fire and forget the cache timeout
		go func(timeout int64) {
			time.Sleep(time.Duration(timeout) * time.Millisecond)
			clearBootstrapCache()
			b.refresh()
		}(config.Cache.LocalTimeout)
	}
	return nil
}

func bootstrapHandler(w http.ResponseWriter, r *http.Request) {
	boot.Scope(r)
	h := w.Header()
	h.Set("Content-Type", "text/javascript")
	e := 60 * 60 * 24
	h.Set("Cache-Control", "max-age="+strconv.Itoa(e)+", must-revalidate")
	bootTmpl.Execute(w, boot.Bootstrap)
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
