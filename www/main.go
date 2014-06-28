package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Bootstrap bootstrap
	Title     string
	Members   []Member
	Groups    []Group
	MapsKey   string
	Version   string
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
	appdata.Title = "Go(lang)MN - Minnesota Go Language Meetup"
	appdata.MapsKey = config.Maps.Key
    appdata.Version = "1"

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		log.Println(pair)
		if pair[0] == "CURRENT_VERSION_ID"{
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


