package members

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
	"members/templates/members.html",))
)


func init() {
	log.Println("initialzing /members")
	http.HandleFunc("/members", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	bootstrap.Bootrapper.Scope(r)
	appdata.Members = bootstrap.Bootrapper.Bootstrap.Members
	err := templates.Execute(w, appdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
