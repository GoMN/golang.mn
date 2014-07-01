package main

import (
	"html/template"
	"net/http"
)

var (
	memberTmpls = template.Must(template.New("layout").Funcs(funcMap).ParseFiles(
	"templates/layout.html",
	"templates/members.html",))
)

func membersHandler(w http.ResponseWriter, r *http.Request) {
	boot.Scope(r)
	appdata.Members = boot.Bootstrap.Members
	err := memberTmpls.Execute(w, appdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
