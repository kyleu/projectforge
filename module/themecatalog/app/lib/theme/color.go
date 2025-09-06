package theme

import (
	"image/color"

	"github.com/muesli/gamut"
)

const dfltNavFG = "#f8f9fa"

func ColorTheme(name string, c color.Color) *Theme {
	light, dark := themeColors(c)
	return &Theme{Key: name, Base: hex(c), Light: light, Dark: dark}
}

func themeColors(c color.Color) (*Colors, *Colors) {
	r, g, b, _ := c.RGBA()
	isLightColor := (r + g + b) > threshold

	var primary, primaryDark color.Color
	if isLightColor {
		primary = gamut.Darker(c, 0.2)
		primaryDark = gamut.Darker(c, 0.7)
	} else {
		primary = gamut.Lighter(c, 0.3)
		primaryDark = c
	}

	primaryTints := gamut.Tints(primary, 6)
	primaryShades := gamut.Shades(primary, 6)
	// complement := gamut.Complementary(primary)
	// complementShades := gamut.Shades(complement, 4)

	darkTints := gamut.Tints(primaryDark, 6)
	darkShades := gamut.Shades(primaryDark, 6)
	// darkComplement := gamut.Complementary(primaryDark)
	// darkComplementTints := gamut.Tints(darkComplement, 4)

	x := Default.Clone("")
	x.Base = hex(c)

	l := x.Light
	l.NavBackground = hex(primaryTints[1])
	l.MenuBackground = hex(primaryTints[4])
	l.MenuSelectedBackground = hex(primaryTints[2])
	l.LinkForeground = hex(darkShades[2])
	l.LinkVisitedForeground = hex(darkShades[2])
	l.BackgroundMuted = hex(primaryTints[5])
	// d.ForegroundMuted = hex(primaryTints[3])
	// l.Border = "1px solid " + hex(primaryTints[2])
	l.Border = "1px solid #cccccc"

	if isLightColor {
		l.NavForeground = "#2a2a2a"
		l.ForegroundMuted = hex(primaryShades[3])
	} else {
		l.NavForeground = dfltNavFG
		l.ForegroundMuted = hex(primaryShades[1])
	}

	d := x.Dark
	d.NavBackground = hex(darkShades[1])
	d.MenuBackground = hex(darkShades[2])
	d.MenuSelectedBackground = hex(darkTints[1])

	// d.LinkForeground = hex(primaryShades[0])
	// d.LinkForeground = "#dddddd"
	d.LinkForeground = hex(primaryTints[4])
	// d.LinkVisitedForeground = hex(primaryShades[1])
	// d.LinkVisitedForeground = "#aaaaaa"
	d.LinkVisitedForeground = hex(primaryTints[1])

	d.BackgroundMuted = hex(darkShades[3])
	// d.ForegroundMuted = hex(darkTints[3])
	// d.Border = "1px solid " + hex(darkTints[1])
	d.Border = "1px solid #444444"

	d.NavForeground = dfltNavFG
	d.ForegroundMuted = hex(darkTints[3])

	return l, d
}

func hex(c color.Color) string {
	return gamut.ToHex(c)
}
