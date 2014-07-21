package bootstrap

import (
	"encoding/json"
	htmltmpl "html/template"
	"net/http"
	"services/conf"
	"strconv"
	"text/template"
	"time"
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
	Bootstrapper.Scope(r)
	h := w.Header()
	h.Set("Content-Type", "text/javascript")
	e := 60 * 60 * 24
	h.Set("Cache-Control", "max-age="+strconv.Itoa(e)+", must-revalidate")
	Bootstrapper.Bootstrap.Version = conf.Config.Version + time.Now().String()
	templates.Execute(w, Bootstrapper.Bootstrap)
}
