package main

import (
	lg "{{{ .Package }}}/app/lib/log"
	"{{{ .Package }}}/app/util"
)

var _rootLogger util.Logger

func main() {
	l, err := lg.InitLogging(true)
	if err != nil {
		println(err)
	}
	_rootLogger = l

	t := util.TimerStart()
	wireFunctions()
	l.Infof("[%s] started in [%s]", util.AppName, t.EndString())
	<-make(chan struct{})
}
