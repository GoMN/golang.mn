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
	Bootstrap bootstrap `json:"bootstrap"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Members   []Member `json:"members"`
	Groups    []Group `json:"groups"`
	MapsKey   string `json:"mapsKey"`
	Year      int `json:"year"`
	Version   string `json:version`
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
	// TODO: replace with proper versioning
	appdata.Version = "1.0.2-" + time.Now().String()
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


