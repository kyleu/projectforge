// Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"github.com/kyleu/projectforge/app/lib/menu"
)

type Breadcrumbs []string

func (b Breadcrumbs) Active(i *menu.Item, path []string) (bool, bool) {
	for idx, x := range path {
		if len(b) <= idx || b[idx] != x {
			return false, false
		}
	}
	if len(i.Children) == 0 {
		return true, true
	}
	if len(path) == len(b) {
		return true, true
	}
	if len(path) < len(b) {
		next := b[len(path)]
		for _, kid := range i.Children {
			if kid.Key == next {
				return true, false
			}
		}
	}
	return true, true
}
