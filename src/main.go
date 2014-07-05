// +build !appengine

package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == ""{
		port = "80"
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
