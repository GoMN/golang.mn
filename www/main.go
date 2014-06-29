package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"
)

type Page struct {
	Bootstrap bootstrap
	Title     string
	Subtitle  string
	Members   []Member
	Groups    []Group
	MapsKey   string
	Version   string
	Year      int
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
	confInit()
	appdata.Title = "Go(lang)MN"
	appdata.Subtitle = "Minnesota Go Language Meetup"
	appdata.MapsKey = config.Maps.Key
	appdata.Version = "1.0.1"
	appdata.Year = time.Now().Year()

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == "CURRENT_VERSION_ID" {
			appdata.Version = pair[1]
			break
		}
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/members", membersHandler)
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("humans.txt", staticHandler)
	http.HandleFunc("robots.txt", staticHandler)
	http.HandleFunc("favicon.ico", staticHandler)
	http.HandleFunc("/static/", staticHandler)
}


