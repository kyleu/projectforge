package task

import (
	"cmp"
	"slices"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Tasks []*Task

func (a Tasks) Sort() {
	slices.SortFunc(a, func(l *Task, r *Task) int {
		return cmp.Compare(l.Key, r.Key)
	})
}

func (a Tasks) Get(key string) *Task {
	return lo.FindOrElse(a, nil, func(x *Task) bool {
		return x.Key == key
	})
}

func (a Tasks) Categories() []string {
	return util.ArraySorted(lo.Uniq(lo.Map(a, func(x *Task, _ int) string {
		return x.Category
	})))
}

func (a Tasks) ByCategory(key string) Tasks {
	return lo.Filter(a, func(x *Task, _ int) bool {
		return x.Category == key
	})
}

func (a Tasks) Keys() []string {
	return lo.Map(a, func(x *Task, _ int) string {
		return x.Key
	})
}

func (a Tasks) Clone() Tasks {
	return lo.Map(a, func(x *Task, _ int) *Task {
		return x.Clone()
	})
}
