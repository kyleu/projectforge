//go:build js
// +build js

package main

import (
	"time"

	"{{{ .Package }}}/app/lib/log"
	"{{{ .Package }}}/app/util"
)

var (
	_rootLogger util.Logger
	_close      chan struct{}
)

func main() {
	logFn := func(level string, occurred time.Time, loggerName string, message string, caller util.ValueMap, stack string, fields util.ValueMap) {
		Log(level, occurred, loggerName, message, caller, stack, fields)
	}
	l, err := log.InitLogging(true, logFn)
	if err != nil {
		println(err)
	}
	_rootLogger = l

	t := util.TimerStart()
	wireFunctions()

	initWASM(l)

	l.Infof("[%s] started in [%s]", util.AppName, t.EndString())
	_close = make(chan struct{})
	<-_close
}
