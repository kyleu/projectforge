package util

import (
	"time"

	"github.com/araddon/dateparse"
	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func TimeTruncate(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	ret := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return &ret
}

func TimeCurrent() time.Time {
	return time.Now()
}

func TimeCurrentP() *time.Time {
	ret := TimeCurrent()
	return &ret
}

func TimeToday() *time.Time {
	return TimeTruncate(TimeCurrentP())
}

func TimeCurrentUnix() int64 {
	return TimeCurrent().Unix()
}

func TimeCurrentMillis() int64 {
	return TimeCurrent().UnixMilli()
}

func TimeCurrentNanos() int64 {
	return TimeCurrent().UnixNano()
}

func TimeRelative(t *time.Time) string {
	if t == nil {
		return "<never>"
	}
	return humanize.Time(*t)
}

func TimeToMap(t time.Time) map[string]any {
	return map[string]any{"epoch": t.UnixMilli(), "iso8601": t.Format("2006-01-02T15:04:05-0700")}
}

func TimeToString(d *time.Time, fmt string) string {
	if d == nil {
		return ""
	}
	return d.Format(fmt)
}

func TimeFromString(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	ret, err := dateparse.ParseLocal(s)
	if err != nil {
		return nil, errors.Errorf("invalid date string [%s]", s)
	}
	return &ret, nil
}

func TimeFromStringSimple(s string) *time.Time {
	ret, _ := TimeFromString(s)
	return ret
}

func TimePlusDays(t *time.Time, days int) *time.Time {
	if t == nil {
		return nil
	}
	ret := t.AddDate(0, 0, days)
	return &ret
}

func TimeMax(ts ...*time.Time) *time.Time {
	return lo.Reduce(ts, func(agg *time.Time, x *time.Time, _ int) *time.Time {
		if x != nil && (agg == nil || x.After(*agg)) {
			return x
		}
		return agg
	}, nil)
}

func TimeMin(ts ...*time.Time) *time.Time {
	return lo.Reduce(ts, func(agg *time.Time, x *time.Time, _ int) *time.Time {
		if x != nil && (agg == nil || x.Before(*agg)) {
			return x
		}
		return agg
	}, nil)
}
