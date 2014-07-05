package about


import (
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
	"about/templates/about.html",))
)

func init(){
	log.Println("initialzing /about")
	http.HandleFunc("/about", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	err := templates.Execute(w, appdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
