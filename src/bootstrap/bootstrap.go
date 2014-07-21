package bootstrap

import (
	"services/logging"
	"net/http"
	"services"
	"services/caching"
	"services/conf"
	"services/meetup"
	"sort"
	"sync"
	"time"
)

var (
	Bootstrapper            = bootstrapper{}
	cache   caching.Cacher
	location, locationErr = time.LoadLocation("America/Chicago")
	meetupSvc             = meetup.NewService()
	log                   = logging.Log{}

)

const BOOTSTRAP_KEY = "bootstrap"

type bootstrapper struct{
	initialized bool
	Bootstrap   bootstrap `json:"bootstrap"`
}

///models
type bootstrap struct{
	Members       []Member `json:"members"`
	MemberCoords  []memberCoord `json:"memberCoords"`
	Calendar      Calendar `json:"calendar"`
	Topics        []meetup.Topic `json:"topics"`
	Version       string `json:"version"`
}
type memberCoord struct {
	Title string `json:"title"`
	Lat   float32 `json:"lat"`
	Lon   float32 `json:"lon"`
}

type Member struct{
	ID       int `json:"id"`
	Joined   int64 `json:"joined"`
	Bio      string `json:"bio"`
	Link     string `json:"link"`
	Name     string `json:"name"`
	City     string `json:"city"`
	State    string `json:"state"`
	Photo    meetup.Photo `json:"photo"`
	Other    meetup.Other `json:"other_services"`
}

func (b *bootstrapper) Scope(r *http.Request) {
	ctx := services.Context{r}
	cache = caching.GetService(ctx)
	log = logging.GetService(ctx)
	meetupSvc.SetContext(ctx)
	if !b.initialized {
		b.initialize()
	}
}

func (b *bootstrapper) Clear() {
	// empty slices to prevent appending latest to cached
	b.Bootstrap.Members = nil
	b.Bootstrap.MemberCoords = nil
	b.Bootstrap.Topics = nil
}

func (b *bootstrapper) initialize() error {
	if b.initialized {
		log.Println("cached bootstrap will be used")
		return nil
	}
	var wg sync.WaitGroup
	var test bootstrap
	var zero = new(bootstrap)
	ok, _ := cache.Get(BOOTSTRAP_KEY, &test)

	if ok && &test != zero {
		log.Println("bootstrap set from cache")
		b.Bootstrap = test
	}else {
		log.Println("bootstrap reinitializing", ok, test)
		b.Clear()
		wg.Add(1)
		go func(boot *bootstrap, svc meetup.MeetupService) {
			defer wg.Done()
			members, err := svc.GetMembers()

			if err != nil {
				log.Printf("ERROR: getting members: %v", err)

			}else {
				var bwg sync.WaitGroup
				bwg.Add(1)
				go func(b *bootstrap) {
					defer bwg.Done()

					for _, m := range members {
						b.Members = append(boot.Members, Member{
							m.ID, m.Joined, m.Bio, m.Link, m.Name, m.City, m.State, m.Photo, m.Other,
						})
						//we want to remove any connection from member and coord exposed publicly
						b.MemberCoords = append(boot.MemberCoords, memberCoord{
							"gopher",
							m.Lat,
							m.Lon,
						})
						b.Topics = append(b.Topics, m.Topics...)
					}
				}(boot)

				bwg.Add(1)
				go func(bs *bootstrap, s meetup.MeetupService) {
					defer bwg.Done()
					calendar, err := b.getMembersCalendar(boot.Members)
					bs.Calendar = calendar
					if err != nil {
						log.Printf("ERROR: building member calendar", err)
					}
				}(boot, svc)

				bwg.Wait()
			}

		}(&b.Bootstrap, meetupSvc)

		//		wg.Add(1)
		//		go func(boot *bootstrap) {
		//			defer wg.Done()
		//			//boot.Calendar =
		//
		//		}(&b.Bootstrap)

		//wait for everything to bootstrap or fail
		wg.Wait()

		//cache this result if any members were returned
		if len(b.Bootstrap.Members) > 0 {
			cache.Set(BOOTSTRAP_KEY, b.Bootstrap)
			b.initialized = true

			//fire and forget the cache timeout
			go func(timeout int64) {
				time.Sleep(time.Duration(timeout) * time.Millisecond)
				clearBootstrapCache()
				b.refresh()
			}(conf.Config.Cache.LocalTimeout)
		}
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
	cache.Set(BOOTSTRAP_KEY, nil)
}

/// event sorting
type ByTimestamp []*meetup.Event

func (a ByTimestamp) Len() int { return len(a) }
func (a ByTimestamp) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }

/// end sorting

type Calendar struct{
	Months []*Month `json:"months"`
	Events []*meetup.Event `json:"events"`
}
type Month struct {
	Name     string `json:"name"`
	StartPos int `json:"startPos"`
	Days     []*Day `json:"days"`
}
type Day struct {
	Date    time.Time `json:"date"`
	Number  int `json:"number"`
	WeekPos int `json:"weekPos"`
	YearDay int
	History bool `json:"history"`
	Events  []meetup.Event `json:"events"`
}

func newCalendar(month time.Month, year int) Calendar {
	c := Calendar{}
	context := time.Date(year, month, 1, 0, 0, 0, 0, location).Local()
	next := context.AddDate(0, 1, 0)

	c.Months = append(c.Months, buildMonth(month, year))
	c.Months = append(c.Months, buildMonth(next.Month(), next.Year()))

	return c
}
func buildMonth(month time.Month, year int) *Month {
	log.Println("building month", month)
	n := time.Now().Local()
	m := Month{}
	m.Name = month.String()
	days := time.Date(year, month+1, 0, 0, 0, 0, 0, location).Day()
	start := time.Date(year, month, 1, 0, 0, 0, 0, location)
	m.StartPos = int(start.Weekday())
	for i := 1; i <= days; i++ {
		d := time.Date(year, month, i, 0, 0, 0, 0, location)
		day := Day{
			d,
			i,
			int(d.Weekday()),
			d.YearDay(),
				n.YearDay() > d.YearDay(),
			[]meetup.Event{},
		}
		m.Days = append(m.Days, &day)
	}
	return &m
}

/// optimize
func (c *Calendar) plotCalendarEvent(e meetup.Event) {
	for _, m := range c.Months {
		for _, d := range m.Days {
			if e.CompareDate == d.Date {
				d.Events = append(d.Events, e)
				break
			}
		}
	}
}
func (b * bootstrapper) getMembersCalendar(members []Member) (Calendar, error) {
	n := time.Now()
	c := newCalendar(n.Month(), n.Year())
	var ids []int
	for _, m := range members {
		ids = append(ids, m.ID)
	}
	groups, err := meetupSvc.GetMemberGroups(ids)

	if err != nil {
		log.Printf("error: retrieving member groups: %v\n", err)
		return c, err
	}

	for _, g := range groups {
		if g.NextEvent.Timestamp > 0 {
			g.NextEvent.Date = time.Unix(0, int64(g.NextEvent.Timestamp)*int64(time.Millisecond)).Local()
			g.NextEvent.YearDay = g.NextEvent.Date.YearDay()
			g.NextEvent.GroupName = g.Name
			g.NextEvent.GroupID = g.ID
			g.NextEvent.GroupURLName = g.URLName
			y, m, d := g.NextEvent.Date.Date();
			g.NextEvent.CompareDate = time.Date(y, m, d, 0, 0, 0, 0, location)

			c.Events = append(c.Events, &g.NextEvent)
			c.plotCalendarEvent(g.NextEvent)
		}
	}
	log.Printf("success: found member groups: %v\n", len(groups))
	sort.Sort(ByTimestamp(c.Events))
	return c, nil
}
