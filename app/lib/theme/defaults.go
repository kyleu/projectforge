// Content managed by Project Forge, see [projectforge.md] for details.
package theme

import (
	"projectforge.dev/projectforge/app/util"
)

var ThemeDefault = func() *Theme {
	nbl := "#50b3b3"
	if o := util.GetEnv("app_nav_color_light"); o != "" {
		nbl = o
	}
	nbd := "#008080"
	if o := util.GetEnv("app_nav_color_dark"); o != "" {
		nbd = o
	}

	return &Theme{
		Key: "default",
		Light: &Colors{
			Border: "1px solid #dddddd", LinkDecoration: "none",
			Foreground: "#000000", ForegroundMuted: "#777777",
			Background: "#ffffff", BackgroundMuted: "#dff0ef",
			LinkForeground: "#346867", LinkVisitedForeground: "#264545",
			NavForeground: "#000000", NavBackground: nbl,
			MenuForeground: "#000000", MenuBackground: "#9cd2d1", MenuSelectedBackground: "#50b3b3", MenuSelectedForeground: "#000000",
			ModalBackdrop: "rgba(77, 77, 77, .7)", Success: "#008000", Error: "#ff0000",
		},
		Dark: &Colors{
			Border: "1px solid #666666", LinkDecoration: "none",
			Foreground: "#ffffff", ForegroundMuted: "#777777",
			Background: "#121212", BackgroundMuted: "#143433",
			LinkForeground: "#7db2b1", LinkVisitedForeground: "#a9cbca",
			NavForeground: "#ffffff", NavBackground: nbd,
			MenuForeground: "#eeeeee", MenuBackground: "#154c4c", MenuSelectedBackground: "#008080", MenuSelectedForeground: "#ffffff",
			ModalBackdrop: "rgba(33, 33, 33, .7)", Success: "#008000", Error: "#ff0000",
		},
	}
}()
