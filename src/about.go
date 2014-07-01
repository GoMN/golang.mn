package main

import (
	"html/template"
	"net/http"
)

var (
	aboutTmpls = template.Must(template.New("layout").Funcs(funcMap).ParseFiles(
	"templates/layout.html",
	"templates/about.html",))
)

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	err := aboutTmpls.Execute(w, appdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
