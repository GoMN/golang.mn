// +build appengine

package http

import(
	"appengine"
	"appengine/urlfetch"
	"net/http"
	"services"
)

var (
	svc = service{}
)

func GetService(c services.Context) *service{
	svc.context = appengine.NewContext(c.Request)
	return &svc
}

type service struct{
	context appengine.Context
}

func(s *service) Get(url string) (*http.Response, error) {
	client := urlfetch.Client(svc.context)
	return client.Get(url)
}
