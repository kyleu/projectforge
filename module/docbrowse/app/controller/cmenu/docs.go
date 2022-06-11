package cmenu

import (
	"context"
	"io/fs"
	"path"
	"strings"

	"golang.org/x/exp/slices"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/doc"
)

var docMenuCached *menu.Item

func docMenu(ctx context.Context, as *app.State, logger util.Logger) *menu.Item {
	if docMenuCached == nil {
		docMenuCached = docMenuCreate(ctx, as, logger)
	}
	return docMenuCached
}

func docMenuCreate(ctx context.Context, as *app.State, logger util.Logger) *menu.Item {
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
		if p == "." || !strings.HasSuffix(p, ".md") {
			continue
		}
		split := strings.Split(p, "/")
		p = strings.TrimSuffix(p, ".md")
		mi := ret
		for idx, comp := range split {
			name := strings.TrimSuffix(comp, ".md")
			addFolder := func() {
				i := &menu.Item{Key: name, Title: util.StringToCamel(name), Icon: "folder"}
				mi.Children = append(mi.Children, i)
				slices.SortFunc(mi.Children, func(l *menu.Item, r *menu.Item) bool {
					return l.Title < r.Title
				})
				mi = i
			}
			if idx == len(split)-1 {
				if strings.HasSuffix(comp, ".md") {
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
		}
	}
	slices.SortFunc(ret.Children, func(l *menu.Item, r *menu.Item) bool {
		return l.Title < r.Title
	})

	return ret
}

func addChild(p string, name string) *menu.Item {
	r := "/" + path.Join("docs", p)
	title := util.StringToCamel(name)
	b, err := doc.FS.ReadFile(p + ".md")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(b), "\n")
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
