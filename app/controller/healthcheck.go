// Content managed by Project Forge, see [projectforge.md] for details.
package controller

import (
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/controller/cutil"
)

func Healthcheck(rc *fasthttp.RequestCtx) {
	x := map[string]string{"status": "OK"}
	_, _ = cutil.RespondJSON(rc, "", x)
}
