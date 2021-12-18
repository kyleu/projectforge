package util

import (
	"time"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

const (
	dateFmtYMD  = "2006-01-02"
	dateFmtFull = "2006-01-02 15:04:05"
	dateFmtJS   = "2006-01-02T15:04:05Z"
)

func TimeToYMD(d *time.Time) string {
	if d == nil {
		return ""
	}
	return d.Format(dateFmtYMD)
}

func TimeFromYMD(s string) (*time.Time, error) {
	return load(s, dateFmtYMD)
}

func TimeToString(d *time.Time) string {
	if d == nil {
		return ""
	}
	return d.Format(dateFmtFull)
}

func TimeFromString(s string) (*time.Time, error) {
	return load(s, dateFmtFull)
}

func TimeFromJS(s string) (*time.Time, error) {
	return load(s, dateFmtJS)
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

func load(s string, f string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	ret, err := time.Parse(f, s)
	if err != nil {
		return nil, errors.New("invalid date string [" + s + "], expected [" + f + "]")
	}
	return &ret, nil
}
