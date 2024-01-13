package controller

import (
	"io"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/project/export/datschema"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vparse"
)

func ParseForm(rc *fasthttp.RequestCtx) {
	Act("parse.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = "Parse Engine"
		return Render(rc, as, &vparse.Form{}, ps)
	})
}

func Parse(rc *fasthttp.RequestCtx) {
	Act("parse", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := &datschema.Schema{}

		frm, err := rc.FormFile("f")
		if err != nil {
			return "", err
		}

		f, err := frm.Open()
		if err != nil {
			return "", err
		}

		b, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}

		err = util.FromJSONStrict(b, &ret)
		if err != nil {
			return "", err
		}

		ps.SetTitleAndData("Parse Result", ret)
		return Render(rc, as, &vparse.Result{Schema: ret}, ps)
	})
}
