// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"time"
)

type Timer struct {
	Started   int64 `json:"started"`
	Completed int64 `json:"complete"`
}

func TimerStart() *Timer {
	return &Timer{Started: TimeCurrentNanos()}
}

func (t *Timer) End() int {
	t.Completed = TimeCurrentNanos()
	return t.Elapsed()
}

func (t *Timer) EndString() string {
	t.End()
	return t.String()
}

func (t *Timer) Elapsed() int {
	if t.Completed == 0 {
		return int((TimeCurrentNanos() - t.Started) / int64(time.Microsecond))
	}
	return int((t.Completed - t.Started) / int64(time.Microsecond))
}

func (t *Timer) String() string {
	return MicrosToMillis(t.Elapsed())
}

type EpochTimestamp struct {
	Epoch   int64      `json:"epoch,omitempty"`
	ISO8601 string     `json:"iso8601,omitempty"`
	Time    *time.Time `json:"-"`
}

func NewEpochTimestamp() *EpochTimestamp {
	t := TimeCurrentP()
	return &EpochTimestamp{Epoch: t.UnixMilli(), ISO8601: t.Format(time.RFC3339), Time: t}
}
