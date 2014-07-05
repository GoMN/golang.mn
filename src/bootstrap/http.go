package bootstrap

import (
	"encoding/json"
	htmltmpl "html/template"
	"text/template"
	"net/http"
	"strconv"
)

var (
	bootFuncMap           = template.FuncMap{
	"marshall": marshall, }
	templates              = template.Must(template.New("bootstrap").Funcs(bootFuncMap).Parse("window.$app.bootstrap = {{ marshall . }};"))
)

func marshall(v interface{}) htmltmpl.JS {
	a, _ := json.Marshal(v)
	return htmltmpl.JS(a)
}

func init() {
	http.HandleFunc("/bootstrap.js", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	Bootrapper.Scope(r)
	h := w.Header()
	h.Set("Content-Type", "text/javascript")
	e := 60 * 60 * 24
	h.Set("Cache-Control", "max-age="+strconv.Itoa(e)+", must-revalidate")
	templates.Execute(w, Bootrapper.Bootstrap)
}
