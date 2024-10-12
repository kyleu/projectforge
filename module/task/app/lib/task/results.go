package task

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Results []*Result

func (r Results) Get(id uuid.UUID) *Result {
	return lo.FindOrElse(r, nil, func(x *Result) bool {
		return x.ID == id
	})
}

func (r Results) ForAct(act string) *Result {
	return lo.FindOrElse(r, nil, func(x *Result) bool {
		return x.Task != nil && x.Task.Key == act
	})
}

func (r Results) Keys() []string {
	return lo.Uniq(lo.Map(r, func(x *Result, _ int) string {
		if x == nil {
			return "-"
		}
		if x.Task == nil {
			return "unknown"
		}
		return x.Task.Key
	}))
}

func (r Results) Statuses() ([]string, map[string]Results) {
	keys := util.ArraySorted(lo.Uniq(lo.Map(r, func(x *Result, _ int) string {
		return x.Status
	})))
	m := make(map[string]Results, len(keys))
	for _, key := range keys {
		m[key] = lo.Filter(r, func(x *Result, _ int) bool {
			return x.Status == key
		})
	}
	return keys, m
}

func (r Results) String() string {
	return fmt.Sprintf("[%s]", strings.Join(lo.Map(r, func(x *Result, _ int) string {
		return x.String()
	}), ", "))
}
