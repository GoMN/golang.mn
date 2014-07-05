// +build !appengine

package http

import (
	"net/http"
	"services"
)

var (
	svc = service{}
)

func GetService(c services.Context) *service{
    return &svc
}

type service struct{

}

func(s *service) Get(url string) (*http.Response, error) {
	return http.Get(url)
}



