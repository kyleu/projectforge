package routes

import (
	"github.com/fasthttp/router"

	"projectforge.dev/projectforge/app/controller/cmodule"
)

func moduleRoutes(r *router.Router) {
	r.GET("/m", cmodule.ModuleList)
	r.GET("/m/{key}", cmodule.ModuleDetail)
	r.GET("/m/{key}/fs", cmodule.ModuleFileRoot)
	r.GET("/m/{key}/fs/{path:*}", cmodule.ModuleFile)
	r.GET("/m/{key}/search", cmodule.ModuleSearch)
}
