package metrics

import (
	"bootstrap"
	"common/models"
	"common/utils"
	"html/template"
	"log"
	"net/http"
)

var (
	appdata = models.AppData()
	templates = template.Must(template.New("layout").Funcs(utils.TemplateFuncMap).ParseFiles(
	"templates/layout.html",
	"metrics/templates/metrics.html",))
)

func init() {
	log.Println("initialzing /metrics")
	http.HandleFunc("/metrics", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	bootstrap.Bootstrapper.Scope(r)
	appdata.Members = bootstrap.Bootstrapper.Bootstrap.Members
	err := templates.Execute(w, appdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
