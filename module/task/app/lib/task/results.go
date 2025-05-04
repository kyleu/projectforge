package task

import (
	"fmt"

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

func (r Results) ForTask(act string) Results {
	return lo.Filter(r, func(x *Result, _ int) bool {
		return x.Task != nil && x.Task.Key == act
	})
}

func (r Results) ForTaskFirst(act string) *Result {
	return lo.FindOrElse(r, nil, func(x *Result) bool {
		return x.Task != nil && x.Task.Key == act
	})
}

func (r Results) ForRun(run string) Results {
	return lo.Filter(r, func(x *Result, _ int) bool {
		return x.Run == run
	})
}

func (r Results) ForRunFirst(run string) *Result {
	return lo.FindOrElse(r, nil, func(x *Result) bool {
		return x.Run == run
	})
}

func (r Results) Keys() []string {
	return lo.Uniq(lo.Map(r, func(x *Result, _ int) string {
		if x == nil {
			return "-"
		}
		if x.Task == nil {
			return util.KeyUnknown
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

func (r Results) AllTags() []string {
	return util.ArraySorted(lo.Uniq(lo.FlatMap(r, func(x *Result, _ int) []string {
		return x.Tags
	})))
}

func (r Results) String() string {
	return fmt.Sprintf("[%s]", util.StringJoin(lo.Map(r, func(x *Result, _ int) string {
		return x.String()
	}), ", "))
}
