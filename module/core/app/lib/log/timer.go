package log

import (
	"fmt"

	"{{{ .Package }}}/app/util"
)

type LogFn func(t *Timer, log string, elapsed int)

// Timer lets you measure laps. It is not safe for concurrent use.
type Timer struct {
	Key         string      `json:"key"`
	Laps        []int       `json:"laps,omitempty"`
	Logs        []string    `json:"logs,omitempty"`
	Log         util.Logger `json:"-"`
	Timer       *util.Timer `json:"-"`
	initial     *util.Timer
	index       int
	acc         int
	persistLogs bool
	fns         []LogFn
}

func NewTimer(key string, persistLogs bool, logger util.Logger, fns ...LogFn) *Timer {
	return &Timer{Key: key, Log: logger, Timer: util.TimerStart(), initial: util.TimerStart(), persistLogs: persistLogs, fns: fns}
}

func (t *Timer) Lap(msg string, args ...any) int {
	t.index++
	out := t.addLog(msg, args...)
	elapsed := t.Timer.End()
	t.acc += elapsed
	t.Laps = append(t.Laps, elapsed)
	t.Timer = util.TimerStart()
	for _, fn := range t.fns {
		fn(t, out, 0)
	}
	return elapsed
}

func (t *Timer) Complete() int {
	_ = t.Lap("completed after [%d] steps in [%s]", t.index, util.MicrosToMillis(t.initial.End()))
	return t.acc
}

func (t *Timer) addLog(msg string, args ...any) string {
	out := fmt.Sprintf("[%s::%d] ", t.Key, t.index) + fmt.Sprintf(msg, args...) + " [" + util.MicrosToMillis(t.Timer.End()) + "]"
	if t.Log != nil {
		t.Log.Infof(out)
	}
	if t.persistLogs {
		t.Logs = append(t.Logs, out)
	}
	return out
}
