package cutil

import (
	"{{{ .Package }}}/app/lib/menu"
)

type Breadcrumbs []string

func (b Breadcrumbs) Active(i *menu.Item, path []string) (bool, bool) {
	for idx, x := range path {
		if len(b) <= idx || b[idx] != x {
			return false, false
		}
	}
	return true, len(i.Children) == 0 || len(path) == len(b)
}
