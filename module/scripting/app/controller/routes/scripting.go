package routes

import (
	"github.com/fasthttp/router"

	"{{{ .Package }}}/app/controller/clib"
)

func scriptingRoutes(r *router.Router) {
	r.GET("/admin/scripting", clib.ScriptingList)
	r.GET("/admin/scripting/new", clib.ScriptingNew)
	r.POST("/admin/scripting/new", clib.ScriptingCreate)
	r.GET("/admin/scripting/{key}", clib.ScriptingDetail)
	r.GET("/admin/scripting/{key}/edit", clib.ScriptingForm)
	r.POST("/admin/scripting/{key}/edit", clib.ScriptingSave)
	r.GET("/admin/scripting/{key}/delete", clib.ScriptingDelete)
}
