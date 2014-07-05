package main

import (
	"log"
	"net/http"
	"services/conf"
	"strconv"
)

func staticHandler(w http.ResponseWriter, r *http.Request) {
	e := 60 * 60 * 24
	w.Header().Set("Cache-Control", "max-age="+strconv.Itoa(e)+", must-revalidate")
	http.ServeFile(w, r, r.URL.Path[1:])
}

func init() {
	log.Println("starting up")
	conf.ConfInit()

	http.HandleFunc("humans.txt", staticHandler)
	http.HandleFunc("robots.txt", staticHandler)
	http.HandleFunc("favicon.ico", staticHandler)
	http.HandleFunc("/static/", staticHandler)
}
