// +build !appengine

package logging

import (
	"log"
	"services"
)

var (
	Printf  = log.Printf
	Println = log.Println
	Fatal   = log.Fatal
)

type Log struct{
	log.Logger
}

func GetService(c services.Context)Log {
	return Log{}
}
