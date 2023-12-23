// Package cmenu - Content managed by Project Forge, see [projectforge.md] for details.
package cmenu

import (
	"cmp"
	"io/fs"
	"path"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/menu"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/doc"
)

var docMenuCached *menu.Item

func docMenu(logger util.Logger) *menu.Item {
	if docMenuCached == nil {
		docMenuCached = docMenuCreate(logger)
	}
	return docMenuCached
}

func docMenuCreate(logger util.Logger) *menu.Item {
	var paths []string
	err := fs.WalkDir(doc.FS, ".", func(path string, d fs.DirEntry, err error) error {
		paths = append(paths, path)
		return err
	})
	if err != nil {
		logger.Errorf("unable to build documentation menu: %+v", err)
	}
	slices.Sort(paths)

	ret := &menu.Item{Key: "docs", Title: "Documentation", Icon: "folder"}
	for _, p := range paths {
		if p == "." || !strings.HasSuffix(p, util.ExtMarkdown) {
			continue
		}
		split := util.StringSplitAndTrim(p, "/")
		p = strings.TrimSuffix(p, util.ExtMarkdown)
		mi := ret
		lo.ForEach(split, func(comp string, idx int) {
			name := strings.TrimSuffix(comp, util.ExtMarkdown)
			addFolder := func() {
				i := &menu.Item{Key: name, Title: util.StringToCamel(name), Icon: "folder"}
				mi.Children = append(mi.Children, i)
				slices.SortFunc(mi.Children, func(l *menu.Item, r *menu.Item) int {
					return cmp.Compare(l.Title, r.Title)
				})
				mi = i
			}
			if idx == len(split)-1 {
				if strings.HasSuffix(comp, util.ExtMarkdown) {
					mi.Children = append(mi.Children, addChild(p, name))
				} else {
					addFolder()
				}
			} else {
				if kid := mi.Children.Get(comp); kid == nil {
					addFolder()
				} else {
					mi = kid
				}
			}
		})
	}
	slices.SortFunc(ret.Children, func(l *menu.Item, r *menu.Item) int {
		return cmp.Compare(l.Title, r.Title)
	})

	return ret
}

func addChild(p string, name string) *menu.Item {
	r := "/" + path.Join("docs", p)
	title := util.StringToCamel(name)
	b, err := doc.FS.ReadFile(p + util.ExtMarkdown)
	if err != nil {
		panic(err)
	}

	lines := util.StringSplitLines(string(b))
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			for strings.HasPrefix(line, "#") {
				line = line[1:]
			}
			title = strings.TrimSpace(line)
			break
		}
	}

	return &menu.Item{Key: name, Title: title, Icon: "file", Route: r}
}
