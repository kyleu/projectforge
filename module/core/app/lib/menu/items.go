package menu

import "github.com/samber/lo"

type Items []*Item

func (i Items) Get(key string) *Item {
	return lo.FindOrElse(i, nil, func(item *Item) bool {
		return item.Key == key
	})
}

func (i Items) Visible() Items {
	return lo.Reject(i, func(item *Item, _ int) bool {
		return item.Hidden
	})
}

func (i Items) GetByPath(path []string) *Item {
	if len(path) == 0 {
		return nil
	}
	ret := i.Get(path[0])
	if ret == nil {
		return nil
	}
	if len(path) > 1 {
		return ret.Children.GetByPath(path[1:])
	}
	return ret
}

func (i Items) Keys() []string {
	return lo.Map(i, func(x *Item, _ int) string {
		return x.Key
	})
}

func (i Items) Titles() []string {
	return lo.Map(i, func(x *Item, _ int) string {
		return x.Title
	})
}
