// Package controller - $PF_IGNORE$
package controller

import (
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views"
)

var homeContent = util.ValueMap{
	"_": util.AppName,
	"urls": map[string]string{
		"TODO": "/todo",
	},
}

func Home(rc *fasthttp.RequestCtx) {
	Act("home", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Data = homeContent
		return Render(rc, as, &views.Home{}, ps)
	})
}
