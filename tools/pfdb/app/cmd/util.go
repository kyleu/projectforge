package cmd

import (
	"sync"

	"projectforge.dev/projectforge/app/lib/log"
	"projectforge.dev/projectforge/app/util"
)

var (
	_initialized = false
	_flags       = &Flags{}
	_logger      util.Logger
)

type Flags struct {
	Debug bool
}

var initMu sync.Mutex

func initIfNeeded() error {
	initMu.Lock()
	defer initMu.Unlock()

	if _initialized {
		return nil
	}
	l, err := log.InitLogging(_flags.Debug)
	if err != nil {
		return err
	}
	util.DEBUG = _flags.Debug
	_logger = l
	_initialized = true
	return nil
}
