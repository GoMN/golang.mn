// +build appengine

package logging

import (
	"appengine"
	"services"
)

type Log struct{
	ctx appengine.Context
}

func GetService(c services.Context)Log {
	l := Log{}
	l.ctx = appengine.NewContext(c.Request)
	return l
}

func (l *Log) Printf(s string, v ...interface{}) {
	l.ctx.Infof(s, v...)
}
func (l *Log) Println(s string, v ...interface{}) {
	l.ctx.Infof(s, v...)
}
func (l *Log) Fatalf(s string, v...interface{}) {
	l.ctx.Criticalf(s, v...)
}
