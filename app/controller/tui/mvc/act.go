package mvc

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

type ActFn func(ts *State, ps *PageState) (Transition, error)

func Act(key string, ts *State, ps *PageState, fn ActFn) (ret Transition, err error) {
	timer := util.TimerStart()
	defer func() {
		if rec := recover(); rec != nil {
			if cast, ok := rec.(error); ok {
				err = errors.Wrapf(cast, "encountered panic in [%s]", key)
			} else {
				err = errors.Wrapf(errors.Errorf("panic in action [%s]", key), "encountered panic of type [%T] in [%s]: %+v", rec, key, rec)
			}
			ps.SetError(err)
			ret = Stay()
		}
		logger := ps.Logger
		if logger == nil && ts != nil {
			logger = ts.Logger
		}
		if logger != nil {
			logger.Debugf("tui act [%s] completed in [%s]", key, util.MicrosToMillis(timer.End()))
		}
	}()
	ret, err = fn(ts, ps)
	if err != nil {
		ps.SetError(err)
	}
	return ret, err
}
