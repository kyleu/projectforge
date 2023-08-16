package routes

import (
	"github.com/fasthttp/router"

	"{{{ .Package }}}/app/controller/clib"
)

func scriptingRoutes(r *router.Router) {
	r.GET("/admin/scripting", clib.ScriptingList)
	r.GET("/admin/scripting/{key}", clib.ScriptingDetail)
}
