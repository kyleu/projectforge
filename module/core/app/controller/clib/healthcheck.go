package clib

import (
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
)

func Healthcheck(rc *fasthttp.RequestCtx) {
	x := util.ValueMap{"status": "OK"}
	_, _ = cutil.RespondJSON(rc, "", x)
}
