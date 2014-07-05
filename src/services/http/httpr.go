package http

import(
	"net/http"
)

type Httpr interface {
    Get(url string) (*http.Response, error)
}
