package theme

import (
	"fmt"
	"image/color"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

const (
	white, black = "#ffffff", "#000000"
	threshold    = (65535 * 3) / 2
)

var Default = func() *Theme {
	nbl := "#83a2b9"
	if o := util.GetEnv("app_nav_color_light"); o != "" {
		nbl = o
	}
	nbd := "#0a2638"
	if o := util.GetEnv("app_nav_color_dark"); o != "" {
		nbd = o
	}

	return &Theme{
		Key: "default",
		Light: &Colors{
			Border: "1px solid #dddddd", LinkDecoration: "none",
			Foreground: "#000000", ForegroundMuted: "#263642",
			Background: "#ffffff", BackgroundMuted: "#e6ecf1",
			LinkForeground: "#0c202e", LinkVisitedForeground: "#0c202e",
			NavForeground: "#2a2a2a", NavBackground: nbl,
			MenuForeground: "#000000", MenuSelectedForeground: "#000000",
			MenuBackground: "#cdd9e3", MenuSelectedBackground: "#9cb4c7",
			ModalBackdrop: "rgba(77, 77, 77, .7)", Success: "#008000", Error: "#ff0000",
		},
		Dark: &Colors{
			Border: "1px solid #666666", LinkDecoration: "none",
			Foreground: "#dddddd", ForegroundMuted: "#93a1af",
			Background: "#121212", BackgroundMuted: "#0d1a24",
			LinkForeground: "#cdd9e3", LinkVisitedForeground: "#83a2b9",
			NavForeground: "#f8f9fa", NavBackground: nbd,
			MenuForeground: "#eeeeee", MenuSelectedForeground: "#dddddd",
			MenuBackground: "#0c202e", MenuSelectedBackground: "#50677d",
			ModalBackdrop: "rgba(33, 33, 33, .7)", Success: "#008000", Error: "#ff0000",
		},
	}
}()

func TextColorFor(clr string) string {
	c, err := ParseHexColor(clr)
	if err != nil {
		return white
	}
	r, g, b, _ := c.RGBA()
	total := r + g + b
	if total < threshold {
		return white
	}
	return black
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
		err = errors.Errorf("invalid length [%d], must be 7 or 4", len(s))
	}
	return ret, err
}
