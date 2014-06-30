package main

import (
	"html/template"
	"net/http"
)

var (
	homeTmpls = template.Must(template.New("layout").Funcs(funcMap).ParseFiles(
	"templates/layout.html",
	"templates/index.html",
	 ))
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	appdata.Members = boot.Bootstrap.Members
	err := homeTmpls.Execute(w, appdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
