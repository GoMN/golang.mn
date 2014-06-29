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
}

type Event struct{
	ID           string `json:"id"`
	Name         string `json:"name"`
	Timestamp    int64 `json:"time"`
	Date         time.Time
	FromNow      int64 `json:"utc_offset"`
}

/// called on module initialization
func init() {
	initialize()
}

func initialize(){
	log.Println("configuring meetup service")
	url_base = config.Meetup.BaseUrl
	meetup_key = config.Meetup.Key
	url_suffix = "sign=true&key="+meetup_key
}

///
/// api below
///

type meetupService struct{
	context appengine.Context
	HttpRequest http.Request
}

func (svc *meetupService) SetContext(c appengine.Context){
	svc.context = c;
}

func (svc * meetupService) getMembers() ([]member, error) {
	if meetup_key == ""{
		initialize()
	}

	url := url_base + "members?group_urlname=golangmn&" + url_suffix

	log.Println("retrieving members", url)

	client := urlfetch.Client(svc.context)

	resp, err :=  client.Get(url)
	//resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	var mr membersResult
	err = json.NewDecoder(resp.Body).Decode(&mr)

	if err != nil {
		log.Fatal("members decode error", err)
	}
	sort.Sort(ByJoined(mr.Results))
	return mr.Results, err
}

func (svc * meetupService) getMember(id int) (*[]byte, error) {
	if meetup_key == ""{
		initialize()
	}
	url := url_base + "member/" + strconv.Itoa(id) + "?" + url_suffix

	client := urlfetch.Client(svc.context)

	resp, err :=  client.Get(url)
	//resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return &body, err
}

func (svc * meetupService) getMemberGroups(ids []int) ([]Group, error) {
	if meetup_key == ""{
		initialize()
	}
	var strids []string
	for _, value := range ids {
		strids = append(strids, strconv.Itoa(value))
	}
	url := url_base + "groups/?member_id=" + strings.Join(strids, ",") + "&fields=next_event&" + url_suffix

	client := urlfetch.Client(svc.context)

	resp, err :=  client.Get(url)
	//resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	var gr groupsResult
	err = json.NewDecoder(resp.Body).Decode(&gr)


	return gr.Results, err
}
