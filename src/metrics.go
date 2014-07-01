package main

import (
	"html/template"
	"net/http"
)

var (
	metricsTmpls = template.Must(template.New("layout").Funcs(funcMap).ParseFiles(
	"templates/layout.html",
	"templates/metrics.html",))
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	boot.Scope(r)
	appdata.Members = boot.Bootstrap.Members
	err := metricsTmpls.Execute(w, appdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
