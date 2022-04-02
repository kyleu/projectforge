package controller

import (
	"context"
	"io/fs"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slices"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/menu"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/doc"
	"{{{ .Package }}}/views/vdoc"
)

func Docs(rc *fasthttp.RequestCtx) {
	act("docs", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		path, _ := RCRequiredString(rc, "path", false)
		if path == "" {
			return "", errors.New("invalid path")
		}

		split := util.StringSplitAndTrim(path, "/")
		bc := []string{"docs"}
		bc = append(bc, split...)

		x, err := doc.HTML(path + ".md")
		if err != nil {
			return "", errors.Wrapf(err, "unable to load documentation from [%s]", path)
		}
		return render(rc, as, &vdoc.MarkdownPage{Title: path, HTML: x}, ps, bc...)
	})
}

func docMenu(ctx context.Context, as *app.State) *menu.Item {
	var paths []string
	err := fs.WalkDir(doc.FS, ".", func(path string, d fs.DirEntry, err error) error {
		paths = append(paths, path)
		return err
	})
	if err != nil {
		as.Logger.Errorf("unable to build documentation menu: %+v", err)
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
				i := &menu.Item{Key: name, Title: name, Icon: "folder"}
				mi.Children = append(mi.Children, i)
				mi = i
			}
			if idx == len(split)-1 {
				if strings.HasSuffix(comp, ".md") {
					r := "/" + path.Join("docs", p)
					mi.Children = append(mi.Children, &menu.Item{Key: name, Title: name, Icon: "file", Route: r})
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

	return ret
}
