// Package log - Content managed by Project Forge, see [projectforge.md] for details.
package log

import (
	"fmt"

	"projectforge.dev/projectforge/app/util"
)

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
}

func NewTimer(key string, persistLogs bool, logger util.Logger) *Timer {
	return &Timer{Key: key, Log: logger, Timer: util.TimerStart(), initial: util.TimerStart(), persistLogs: persistLogs}
}

func (l *Timer) Lap(msg string, args ...any) int {
	l.index++
	l.addLog(msg, args...)
	elapsed := l.Timer.End()
	l.acc += elapsed
	l.Laps = append(l.Laps, elapsed)
	l.Timer = util.TimerStart()
	return elapsed
}

func (l *Timer) Complete() int {
	msg := fmt.Sprintf("completed after [%d] steps in [%s]", l.index, util.MicrosToMillis(l.initial.End()))
	out := fmt.Sprintf("[%s::%d] ", l.Key, l.index) + msg + " [" + util.MicrosToMillis(l.Timer.End()) + "]"
	l.addLog(out)
	return l.acc
}

func (l *Timer) addLog(msg string, args ...any) {
	out := fmt.Sprintf("[%s::%d] ", l.Key, l.index) + fmt.Sprintf(msg, args...) + " [" + util.MicrosToMillis(l.Timer.End()) + "]"
	if l.Log == nil {
		fmt.Println(out) //nolint:forbidigo
	} else {
		l.Log.Infof(out)
	}
	if l.persistLogs {
		l.Logs = append(l.Logs, out)
	}
}
