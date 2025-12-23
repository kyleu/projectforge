package cmenu

import (
	"cmp"
	"io/fs"
	"slices"

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
	for _, pth := range util.ArraySorted(paths) {
		p := util.Str(pth)

		if p.Equal(".") || p.HasPrefix("module/") || !p.HasSuffix(util.ExtMarkdown) {
			continue
		}
		split := p.SplitAndTrim("/")
		p = p.TrimSuffix(util.ExtMarkdown)
		mi := ret
		lo.ForEach(split, func(comp util.Str, idx int) {
			name := comp.TrimSuffix(util.ExtMarkdown)
			addFolder := func() {
				i := &menu.Item{Key: name.String(), Title: comp.ToProper().String(), Icon: "folder"}
				mi.Children = append(mi.Children, i)
				slices.SortFunc(mi.Children, func(l *menu.Item, r *menu.Item) int {
					return cmp.Compare(l.Title, r.Title)
				})
				mi = i
			}
			if idx == len(split)-1 {
				if comp.HasSuffix(util.ExtMarkdown) {
					mi.Children = append(mi.Children, addChild(p.String(), name.String()))
				} else {
					addFolder()
				}
			} else {
				if kid := mi.Children.Get(comp.String()); kid == nil {
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
	for _, l := range lines {
		line := util.Str(l)
		if line.HasPrefix("#") {
			for line.HasPrefix("#") {
				line = line[1:]
			}
			title = line.TrimSpace().String()
			break
		}
	}

	return &menu.Item{Key: name, Title: title, Icon: "file", Route: r}
}
