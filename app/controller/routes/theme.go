// Package routes - Content managed by Project Forge, see [projectforge.md] for details.
package routes

import (
	"github.com/fasthttp/router"

	"projectforge.dev/projectforge/app/controller/clib"
)

func themeRoutes(r *router.Router) {
	r.GET("/theme", clib.ThemeList)
	r.GET("/theme/{key}", clib.ThemeEdit)
	r.POST("/theme/{key}", clib.ThemeSave)
	r.GET("/theme/{key}/remove", clib.ThemeRemove)

	r.GET("/theme/color/{color}", clib.ThemeColor)
	r.GET("/theme/color/edit", clib.ThemeColorEdit)
	r.GET("/theme/palette/{palette}", clib.ThemePalette)
	r.GET("/theme/palette/{palette}/{theme}", clib.ThemePaletteEdit)
}
