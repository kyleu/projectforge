// $PF_GENERATE_ONCE$
//go:build js
// +build js

package wasm

import (
	"time"

	"{{{ .Package }}}/app/lib/log"
	"{{{ .Package }}}/app/util"
)

type WASM struct {
	logger  util.Logger
	closeCh chan struct{}
}

func NewWASM() (*WASM, error) {
	ret := &WASM{CloseCh: make(chan struct{})}

	logFn := func(level string, occurred time.Time, loggerName string, message string, caller util.ValueMap, stack string, fields util.ValueMap) {
		ret.Log(level, occurred, loggerName, message, caller, stack, fields)
	}
	l, err := log.InitLogging(true, logFn)
	if err != nil {
		return nil, err
	}
	ret.logger = l

	t := util.TimerStart()
	ret.wireFunctions()
	l.Infof("[%s] started in [%s]", util.AppName, t.EndString())

	return ret, nil
}

func (w *WASM) Run() error {
	<-w.closeCh
	return nil
}
