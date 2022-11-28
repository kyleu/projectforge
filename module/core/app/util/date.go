package util

import (
	"time"

	"github.com/araddon/dateparse"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

const (
	dateFmtYMD  = "2006-01-02"
	dateFmtFull = "2006-01-02 15:04:05"
	dateFmtHTML = "2006-01-02T15:04:05"
	dateFmtJS   = "2006-01-02T15:04:05Z"
)

func TimeToday() *time.Time {
	t := time.Now()
	return TimeTruncate(&t)
}

func TimeTruncate(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	ret := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return &ret
}

func TimeCurrentMillis() int64 {
	return time.Now().UnixMilli()
}

func TimeRelative(t *time.Time) string {
	if t == nil {
		return "<never>"
	}
	return humanize.Time(*t)
}

func TimeToString(d *time.Time, fmt string) string {
	if d == nil {
		return ""
	}
	return d.Format(fmt)
}

func TimeToYMD(d *time.Time) string {
	return TimeToString(d, dateFmtYMD)
}

func TimeToFull(d *time.Time) string {
	return TimeToString(d, dateFmtFull)
}

func TimeToHTML(d *time.Time) string {
	return TimeToString(d, dateFmtHTML)
}

func TimeToJS(d *time.Time) string {
	return TimeToString(d, dateFmtJS)
}

func TimeFromString(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	ret, err := dateparse.ParseLocal(s)
	if err != nil {
		return nil, errors.New("invalid date string [" + s + "]")
	}
	return &ret, nil
}

func TimeFromStringFmt(s string, fmt string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	ret, err := time.Parse(fmt, s)
	if err != nil {
		return nil, errors.New("invalid date string [" + s + "], expected [" + fmt + "]")
	}
	return &ret, nil
}

func TimeFromYMD(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtYMD)
}

func TimeFromFull(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtFull)
}

func TimeFromHTML(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtHTML)
}

func TimeFromJS(s string) (*time.Time, error) {
	return TimeFromStringFmt(s, dateFmtJS)
}

func TimeToMap(t time.Time) map[string]any {
	return map[string]any{"epoch": t.UnixMilli(), "iso8601": t.Format("2006-01-02T15:04:05-0700")}
}

func NowPointer() *time.Time {
	ret := time.Now()
	return &ret
}
