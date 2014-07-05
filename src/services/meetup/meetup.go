package meetup

import (
	"sync"
	"io/ioutil"
	"encoding/json"
	"log"
	"services"
	"services/conf"
	"services/http"
	"sort"
	"strings"
	"strconv"
	"time"
)

var (
	httpSvc http.Httpr
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
	CompareDate  time.Time
	FromNow      int64 `json:"utc_offset"`
	YearDay      int
}

/// called on module initialization
func init() {
	initialize()
	log.Println("location", location)
}

func initialize() {
	url_base = conf.Config.Meetup.BaseUrl
	meetup_key = conf.Config.Meetup.Key
	url_suffix = "sign=true&key="+meetup_key
}

func NewService() MeetupService{
	return MeetupService{}
}

///
/// api below
///

type MeetupService struct{
	context     services.Context
}

func (svc *MeetupService) SetContext(c services.Context) {
	svc.context = c;
	httpSvc = http.GetService(c)
}

func (svc * MeetupService) GetMembers() ([]member, error) {
	if meetup_key == "" {
		initialize()
	}

	url := url_base + "members?group_urlname=golangmn&" + url_suffix

	resp, err := httpSvc.Get(url)

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

func (svc * MeetupService) getMember(id int) (*[]byte, error) {
	if meetup_key == "" {
		initialize()
	}
	url := url_base + "member/" + strconv.Itoa(id) + "?" + url_suffix

	resp, err := httpSvc.Get(url)

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

func (svc * MeetupService) GetMemberGroups(ids []int) ([]Group, error) {
	totalGroups := groupsResult{}
	var err error
	max := 40
	var wg sync.WaitGroup
	if meetup_key == "" {
		initialize()
	}
	var strids []string
	for _, value := range ids {
		strids = append(strids, strconv.Itoa(value))
	}

	l := len(strids) / max

	// we search for member groups by passing in an array of member ids
	// url length limitation requires we batch these request at around 40 members
	// we use a wait group to send the batches asynchronously
	for i := 0; i < l; i +=1 {
		wg.Add(1)
		num := (i * max) + max
		if (num > (len(strids)-1)) {
			num = len(strids)-1
		}
		cids := strids[max*i:num]
		go func(g *groupsResult, ids []string) {
			defer wg.Done()
			url := url_base + "groups/?member_id=" + strings.Join(ids, ",") + "&fields=next_event&" + url_suffix

			resp, err := httpSvc.Get(url)
			var gr groupsResult
			err = json.NewDecoder(resp.Body).Decode(&gr)
			if err != nil {
				log.Println("error retrieving groups", err)
			}else {
				g.Results = append(g.Results, gr.Results...)
			}
		}(&totalGroups, cids)

	}
	wg.Wait()

	return totalGroups.Results, err
}
