package theme

import (
	"fmt"
	"image/color"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
)

const (
	white, black = "#ffffff", "#000000"
	threshold    = (65535 * 3) / 2
)

var Default = func() *Theme {
	nbl := "{{{ .Theme.Light.NavBackground }}}"
	if o := util.GetEnv("app_nav_color_light"); o != "" {
		nbl = o
	}
	nbd := "{{{ .Theme.Dark.NavBackground }}}"
	if o := util.GetEnv("app_nav_color_dark"); o != "" {
		nbd = o
	}

	return &Theme{
		Key: "default",
		Light: &Colors{
			Border: "{{{ .Theme.Light.Border }}}", LinkDecoration: "{{{ .Theme.Light.LinkDecoration }}}",
			Foreground: "{{{ .Theme.Light.Foreground }}}", ForegroundMuted: "{{{ .Theme.Light.ForegroundMuted }}}",
			Background: "{{{ .Theme.Light.Background }}}", BackgroundMuted: "{{{ .Theme.Light.BackgroundMuted }}}",
			LinkForeground: "{{{ .Theme.Light.LinkForeground }}}", LinkVisitedForeground: "{{{ .Theme.Light.LinkVisitedForeground }}}",
			NavForeground: "{{{ .Theme.Light.NavForeground }}}", NavBackground: nbl,
			MenuForeground: "{{{ .Theme.Light.MenuForeground }}}", MenuSelectedForeground: "{{{ .Theme.Light.MenuSelectedForeground }}}",
			MenuBackground: "{{{ .Theme.Light.MenuBackground }}}", MenuSelectedBackground: "{{{ .Theme.Light.MenuSelectedBackground }}}",
			ModalBackdrop: "{{{ .Theme.Light.ModalBackdrop }}}", Success: "{{{ .Theme.Light.Success }}}", Error: "{{{ .Theme.Light.Error }}}",
		},
		Dark: &Colors{
			Border: "{{{ .Theme.Dark.Border }}}", LinkDecoration: "{{{ .Theme.Dark.LinkDecoration }}}",
			Foreground: "{{{ .Theme.Dark.Foreground }}}", ForegroundMuted: "{{{ .Theme.Dark.ForegroundMuted }}}",
			Background: "{{{ .Theme.Dark.Background }}}", BackgroundMuted: "{{{ .Theme.Dark.BackgroundMuted }}}",
			LinkForeground: "{{{ .Theme.Dark.LinkForeground }}}", LinkVisitedForeground: "{{{ .Theme.Dark.LinkVisitedForeground }}}",
			NavForeground: "{{{ .Theme.Dark.NavForeground }}}", NavBackground: nbd,
			MenuForeground: "{{{ .Theme.Dark.MenuForeground }}}", MenuSelectedForeground: "{{{ .Theme.Dark.MenuSelectedForeground }}}",
			MenuBackground: "{{{ .Theme.Dark.MenuBackground }}}", MenuSelectedBackground: "{{{ .Theme.Dark.MenuSelectedBackground }}}",
			ModalBackdrop: "{{{ .Theme.Dark.ModalBackdrop }}}", Success: "{{{ .Theme.Dark.Success }}}", Error: "{{{ .Theme.Dark.Error }}}",
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
