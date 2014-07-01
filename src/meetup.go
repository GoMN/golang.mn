package main

import (
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"strconv"
	"time"
)

var (
	location, locationErr = time.LoadLocation("America/Chicago")
	url_base string
	meetup_key string
	url_suffix string
)

/// models
type membersResult struct {
	Results []member `json:"results"`
}
type groupsResult struct {
	Results []Group `json:"results"`
}

/// keep primary object private
type member struct{
	ID       int `json:"id"`
	Joined   int64 `json:"joined"`
	Bio      string `json:"bio"`
	Lat      float32 `json:"lat"`
	Lon      float32 `json:"lon"`
	Link     string `json:"link"`
	Name     string `json:"name"`
	City     string `json:"city"`
	State    string `json:"state"`
	Photo    Photo `json:"photo"`
	Topics   []Topic `json:"topics"`
	Other    Other `json:"other_services"`
}

type Member struct{
	ID       int `json:"id"`
	Joined   int64 `json:"joined"`
	Bio      string `json:"bio"`
	Link     string `json:"link"`
	Name     string `json:"name"`
	City     string `json:"city"`
	State    string `json:"state"`
	Photo    Photo `json:"photo"`
	Other    Other `json:"other_services"`
}

/// member sorting
type ByJoined []member

func (a ByJoined) Len() int { return len(a) }
func (a ByJoined) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByJoined) Less(i, j int) bool { return a[i].Joined < a[j].Joined }

/// end sorting

type Other struct {
	Flickr   Social `json:"flickr"`
	Tumblr   Social `json:"tumblr"`
	Twitter  Social `json:"twitter"`
	Linkedin Social `json:"linkedin"`
	Facebook SocialInt `json:"facebook"`
}

/// json eoncoder could not decode a json int into a string...whoa?
type SocialInt struct {
	Identifier int `json:"identifier"`
}

type Social struct {
	Identifier string `json:"identifier"`
}

/// end whoa

type Photo struct{
	ID    int `json:"photo_id"`
	Thumb string `json:"thumb_link"`
	Full  string `json:"photo_link"`
}

type Topic struct {
	ID     int `json:"id"`
	Name   string  `json:"name"`
	URLKey string  `json:"urlkey"`
}

type Group struct {
	ID        int `json:"id"`
	Name      string `json:"name"`
	NextEvent Event `json:"next_event"`
	URLName   string `json:"urlname"`
}

type Event struct{
	ID           string `json:"id"`
	Name         string `json:"name"`
	GroupName    string `json:"groupName"`
	GroupURLName string `json:"groupURLName"`
	GroupID      int `json:"groupID"`
	Timestamp    int `json:"time"`
	Date         time.Time `json:"date"`
	compareDate  time.Time
	FromNow      int64 `json:"utc_offset"`
	YearDay      int
}

/// event sorting
type ByTimestamp []*Event

func (a ByTimestamp) Len() int { return len(a) }
func (a ByTimestamp) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }

/// end sorting

type Calendar struct{
	Months []*Month `json:"months"`
	Events []*Event `json:"events"`
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
	Events  []Event `json:"events"`
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
			[]Event{},
		}
		m.Days = append(m.Days, &day)
	}
	return &m
}

/// optimize
func (c *Calendar) plotCalendarEvent(e Event) {
	for _, m := range c.Months {
		for _, d := range m.Days {
			if e.compareDate == d.Date {
				d.Events = append(d.Events, e)
				break
			}
		}
	}
}
func (svc * meetupService) getMembersCalendar(members []Member) (Calendar, error) {
	n := time.Now()
	c := newCalendar(n.Month(), n.Year())
	var ids []int
	for _, m := range members {
		ids = append(ids, m.ID)
	}
	groups, err := svc.getMemberGroups(ids)

	if err != nil{
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
			g.NextEvent.compareDate = time.Date(y, m, d, 0, 0, 0, 0, location)

			c.Events = append(c.Events, &g.NextEvent)
			c.plotCalendarEvent(g.NextEvent)
		}
	}
	log.Printf("success: found member groups: %v\n", len(groups))
	sort.Sort(ByTimestamp(c.Events))
	return c, nil
}

/// called on module initialization
func init() {
	initialize()
	log.Println("location", location)
}

func initialize() {
	url_base = config.Meetup.BaseUrl
	meetup_key = config.Meetup.Key
	url_suffix = "sign=true&key="+meetup_key
}

///
/// api below
///

type meetupService struct{
	context     appengine.Context
	HttpRequest http.Request
}

func (svc *meetupService) SetContext(c appengine.Context) {
	svc.context = c;
}

func (svc * meetupService) getMembers() ([]member, error) {
	if meetup_key == "" {
		initialize()
	}

	url := url_base + "members?group_urlname=golangmn&" + url_suffix

	client := urlfetch.Client(svc.context)

	resp, err := client.Get(url)

	if err != nil {
		return []member{}, err
	}

	var mr membersResult
	err = json.NewDecoder(resp.Body).Decode(&mr)

	if err != nil {
		log.Printf("members decode error: %v", err)
		return []member{}, err
	}
	sort.Sort(ByJoined(mr.Results))
	return mr.Results, err
}

func (svc * meetupService) getMember(id int) (*[]byte, error) {
	if meetup_key == "" {
		initialize()
	}
	url := url_base + "member/" + strconv.Itoa(id) + "?" + url_suffix

	client := urlfetch.Client(svc.context)

	resp, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return &body, err
	}

	return &body, err
}

func (svc * meetupService) getMemberGroups(ids []int) ([]Group, error) {
	if meetup_key == "" {
		initialize()
	}
	var strids []string
	for _, value := range ids {
		strids = append(strids, strconv.Itoa(value))
	}
	url := url_base + "groups/?member_id=" + strings.Join(strids, ",") + "&fields=next_event&" + url_suffix

	client := urlfetch.Client(svc.context)

	resp, err := client.Get(url)

	if err != nil {
		return []Group{}, err
	}

	var gr groupsResult
	err = json.NewDecoder(resp.Body).Decode(&gr)


	return gr.Results, err
}
