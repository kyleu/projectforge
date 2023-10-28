package controller

import (
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views"
)

func Testbed(rc *fasthttp.RequestCtx) {
	Act("testbed", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := "Testbed!"
		frm, _ := cutil.ParseForm(rc)
		if len(frm) > 0 {
			ret = frm.GetStringOpt("x")
		}
		ps.SetTitleAndData("Testbed", ret)
		return Render(rc, as, &views.Testbed{Param: ret}, ps)
	})
}
