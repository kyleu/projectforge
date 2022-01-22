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
	return true, len(i.Children) == 0 || len(path) == len(b)
}
