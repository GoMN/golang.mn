package main

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type Page struct {
	Bootstrap bootstrap
	Title     string
	Members   []Member
	Groups    []Group
	MapsKey   string
}

func addInt(num int, num2 int) int {
	return num + num2
}

func marshall(v interface{}) template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}

var (
	appdata = Page{}
	funcMap = template.FuncMap{
	"addint": addInt,
	"marshall": marshall, }
)

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func init() {
	appdata.Title = "Go(lang)MN - Minnesota Go Language Meetup"
	appdata.MapsKey = config.Maps.Key

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/members", membersHandler)
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("humans.txt", staticHandler)
	http.HandleFunc("robots.txt", staticHandler)
	http.HandleFunc("favicon.ico", staticHandler)
	http.HandleFunc("/static/", staticHandler)
}


