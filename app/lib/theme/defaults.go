// Content managed by Project Forge, see [projectforge.md] for details.
package theme

import (
	"projectforge.dev/projectforge/app/util"
)

var ThemeDefault = func() *Theme {
	nbl := "#4f9abd"
	if o := util.GetEnv("app_nav_color_light"); o != "" {
		nbl = o
	}
	nbd := "#2d414e"
	if o := util.GetEnv("app_nav_color_dark"); o != "" {
		nbd = o
	}

	return &Theme{
		Key: "default",
		Light: &Colors{
			Border: "1px solid #dddddd", LinkDecoration: "none",
			Foreground: "#000000", ForegroundMuted: "#999999",
			Background: "#ffffff", BackgroundMuted: "#eeeeee",
			LinkForeground: "#2d414e", LinkVisitedForeground: "#406379",
			NavForeground: "#000000", NavBackground: nbl,
			MenuForeground: "#000000", MenuBackground: "#f0f8ff", MenuSelectedBackground: "#faebd7", MenuSelectedForeground: "#000000",
			ModalBackdrop: "rgba(77, 77, 77, .7)", Success: "#008000", Error: "#FF0000",
		},
		Dark: &Colors{
			Border: "1px solid #666666", LinkDecoration: "none",
			Foreground: "#ffffff", ForegroundMuted: "#999999",
			Background: "#121212", BackgroundMuted: "#333333",
			LinkForeground: "#2d414e", LinkVisitedForeground: "#aaaaaa",
			NavForeground: "#ffffff", NavBackground: nbd,
			MenuForeground: "#eeeeee", MenuBackground: "#171f24", MenuSelectedBackground: "#333333", MenuSelectedForeground: "#ffffff",
			ModalBackdrop: "rgba(33, 33, 33, .7)", Success: "#008000", Error: "#FF0000",
		},
	}
}()
