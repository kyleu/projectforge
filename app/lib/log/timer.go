// Package log - Content managed by Project Forge, see [projectforge.md] for details.
package log

import (
	"fmt"

	"projectforge.dev/projectforge/app/util"
)

// Timer lets you measure laps. It is not safe for concurrent use.
type Timer struct {
	Key     string      `json:"key"`
	Log     util.Logger `json:"-"`
	Timer   *util.Timer `json:"-"`
	initial *util.Timer
	index   int
	acc     int
}

func NewTimer(key string, logger util.Logger) *Timer {
	return &Timer{Key: key, Log: logger, Timer: util.TimerStart(), initial: util.TimerStart()}
}

func (l *Timer) Lap(msg string, args ...any) {
	l.index++
	out := fmt.Sprintf("[%s::%d] ", l.Key, l.index) + fmt.Sprintf(msg, args...) + " [" + util.MicrosToMillis(l.Timer.End()) + "]"
	if l.Log == nil {
		fmt.Println(out) //nolint:forbidigo
	} else {
		l.Log.Infof(out)
	}
	l.acc += l.Timer.End()
	l.Timer = util.TimerStart()
}

func (l *Timer) Complete() int {
	msg := fmt.Sprintf("completed after [%d] steps in [%s]", l.index, util.MicrosToMillis(l.initial.End()))
	out := fmt.Sprintf("[%s::%d] ", l.Key, l.index) + msg + " [" + util.MicrosToMillis(l.Timer.End()) + "]"
	if l.Log == nil {
		fmt.Println(out) //nolint:forbidigo
	} else {
		l.Log.Infof(out)
	}
	l.acc += l.Timer.End()
	return l.acc
}
