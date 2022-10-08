package routes

import (
	"github.com/fasthttp/router"

	"{{{ .Package }}}/app/controller/clib"
)

func execRoutes(r *router.Router) {
	r.GET("/admin/exec", clib.ExecList)
	r.GET("/admin/exec/new", clib.ExecForm)
	r.POST("/admin/exec/new", clib.ExecNew)
	r.GET("/admin/exec/{key}/{idx}", clib.ExecDetail){{{ if .HasModule "websocket" }}}
	r.GET("/admin/exec/{key}/{idx}/connect", clib.ExecSocket){{{ end }}}
	r.GET("/admin/exec/{key}/{idx}/kill", clib.ExecKill)
}
