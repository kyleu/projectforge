package clib

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/doc"
	"{{{ .Package }}}/views/vdoc"
)

func Docs(rc *fasthttp.RequestCtx) {
	controller.Act("docs", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		pth, _ := cutil.RCRequiredString(rc, "path", false)
		if pth == "" {
			return "", errors.New("invalid path")
		}

		bc := []string{"docs"}
		bc = append(bc, util.StringSplitAndTrim(pth, "/")...)

		title, x, err := doc.HTML("doc:"+pth, pth+".md", func(s string) (string, string, error) {
			return cutil.FormatCleanMarkup(s, "file")
		})
		if err != nil {
			return "", errors.Wrapf(err, "unable to load documentation from [%s]", pth)
		}
		ps.Title = title
		ps.Data, _ = doc.Content(pth + ".md")
		return controller.Render(rc, as, &vdoc.MarkdownPage{Title: pth, HTML: x}, ps, bc...)
	})
}
