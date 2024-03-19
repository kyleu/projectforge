package controller

import (
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/views"
)

func Testbed(w http.ResponseWriter, r *http.Request) {
	Act("testbed", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := "Testbed!"
		frm, _ := cutil.ParseForm(r, ps.RequestBody)
		if len(frm) > 0 {
			ret = frm.GetStringOpt("x")
		}
		ps.SetTitleAndData("Testbed", ret)
		return Render(w, r, as, &views.Testbed{Param: ret}, ps)
	})
}
