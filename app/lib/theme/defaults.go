// Package theme - Content managed by Project Forge, see [projectforge.md] for details.
package theme

import (
	"fmt"
	"image/color"

	"projectforge.dev/projectforge/app/util"
)

const threshold = (65535 * 3) / 2

var Default = func() *Theme {
	nbl := "#72a0c1"
	if o := util.GetEnv("app_nav_color_light"); o != "" {
		nbl = o
	}
	nbd := "#00415d"
	if o := util.GetEnv("app_nav_color_dark"); o != "" {
		nbd = o
	}

	return &Theme{
		Key: "default",
		Light: &Colors{
			Border: "1px solid #dddddd", LinkDecoration: "none",
			Foreground: "#000000", ForegroundMuted: "#777777",
			Background: "#ffffff", BackgroundMuted: "#e3ebf2",
			LinkForeground: "#455d6f", LinkVisitedForeground: "#303e4a",
			NavForeground: "#000000", NavBackground: nbl,
			MenuForeground: "#000000", MenuSelectedForeground: "#000000",
			MenuBackground: "#acc5da", MenuSelectedBackground: "#72a0c1",
			ModalBackdrop: "rgba(77, 77, 77, .7)", Success: "#008000", Error: "#ff0000",
		},
		Dark: &Colors{
			Border: "1px solid #666666", LinkDecoration: "none",
			Foreground: "#ffffff", ForegroundMuted: "#777777",
			Background: "#121212", BackgroundMuted: "#1f4055",
			LinkForeground: "#6f889b", LinkVisitedForeground: "#9eaebb",
			NavForeground: "#ffffff", NavBackground: nbd,
			MenuForeground: "#eeeeee", MenuSelectedForeground: "#ffffff",
			MenuBackground: "#0e2939", MenuSelectedBackground: "#00415d",
			ModalBackdrop: "rgba(33, 33, 33, .7)", Success: "#008000", Error: "#ff0000",
		},
	}
}()

func TextColorFor(clr string) string {
	c, err := ParseHexColor(clr)
	if err != nil {
		return "#ffffff"
	}
	r, g, b, _ := c.RGBA()
	total := r + g + b
	if total < threshold {
		return "#ffffff"
	}
	return "#000000"
}

func ParseHexColor(s string) (color.RGBA, error) {
	ret := color.RGBA{A: 0xff}
	var err error
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &ret.R, &ret.G, &ret.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &ret.R, &ret.G, &ret.B)
		// Double the hex digits:
		ret.R *= 17
		ret.G *= 17
		ret.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}
	return ret, err
}
