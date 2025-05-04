package cmenu

import (
	"cmp"
	"io/fs"
	"slices"
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/doc"
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
	err := fs.WalkDir(doc.FS, ".", func(path string, _ fs.DirEntry, err error) error {
		paths = append(paths, path)
		return err
	})
	if err != nil {
		logger.Errorf("unable to build documentation menu: %+v", err)
	}

	ret := &menu.Item{Key: "docs", Title: "Documentation", Icon: "folder"}
	for _, p := range util.ArraySorted(paths) {
		if p == "." || strings.HasPrefix(p, "module/") || !strings.HasSuffix(p, util.ExtMarkdown) {
			continue
		}
		split := util.StringSplitAndTrim(p, "/")
		p = strings.TrimSuffix(p, util.ExtMarkdown)
		mi := ret
		lo.ForEach(split, func(comp string, idx int) {
			name := strings.TrimSuffix(comp, util.ExtMarkdown)
			addFolder := func() {
				i := &menu.Item{Key: name, Title: util.StringToProper(name), Icon: "folder"}
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
	r := "/" + util.StringPath("docs", p)
	title := util.StringToProper(name)
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
