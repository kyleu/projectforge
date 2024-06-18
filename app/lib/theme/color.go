package theme

import (
	"image/color"

	"github.com/muesli/gamut"
)

func ColorTheme(name string, c color.Color) *Theme {
	light, dark := themeColors(c)
	return &Theme{Key: name, Light: light, Dark: dark}
}

func themeColors(c color.Color) (*Colors, *Colors) {
	light, dark := c, c

	r, g, b, _ := c.RGBA()
	if total := r + g + b; total < threshold {
		light = gamut.Lighter(c, 0.4)
	} else {
		dark = gamut.Darker(c, 0.6)
	}

	lightTints := gamut.Tints(light, 4)
	lightShades := gamut.Shades(light, 4)

	darkTints := gamut.Tints(dark, 4)
	darkShades := gamut.Shades(dark, 4)

	x := Default.Clone("")

	l := x.Light
	l.NavBackground = hex(light)
	l.MenuBackground = hex(lightTints[1])
	l.MenuSelectedBackground = hex(light)
	l.LinkForeground = hex(lightShades[1])
	l.LinkVisitedForeground = hex(lightShades[2])
	l.BackgroundMuted = hex(lightTints[3])
	// l.ForegroundMuted = hex(lightShades[3])

	d := x.Dark
	d.NavBackground = hex(dark)
	d.MenuBackground = hex(darkShades[1])
	d.MenuSelectedBackground = hex(dark)
	d.LinkForeground = hex(darkTints[1])
	d.LinkVisitedForeground = hex(darkTints[2])
	d.BackgroundMuted = hex(darkShades[1])
	// d.ForegroundMuted = hex(darkTints[3])

	return l, d
}

func hex(c color.Color) string {
	return gamut.ToHex(c)
}
