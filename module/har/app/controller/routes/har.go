package routes

import (
	"github.com/fasthttp/router"

	"{{{ .Package }}}/app/controller/clib"
)

func harRoutes(r *router.Router) {
	r.GET("/har", clib.HarList)
	r.POST("/har", clib.HarUpload)
	r.GET("/har/{key}", clib.HarDetail)
	r.GET("/har/{key}/delete", clib.HarDelete)
	r.GET("/har/{key}/trim", clib.HarTrim)
}
