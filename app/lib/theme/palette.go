// Package theme - Content managed by Project Forge, see [projectforge.md] for details.
package theme

import (
	"cmp"
	"slices"

	"github.com/muesli/gamut"
	"github.com/muesli/gamut/palette"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

func PaletteTheme(pal string, key string) (*Theme, error) {
	p, err := getPalette(pal)
	if err != nil {
		return nil, err
	}
	matched := p.Filter(key)
	if len(matched) == 0 {
		return nil, errors.Errorf("no color available in palette [%s] with key [%s]", pal, key)
	}
	return ColorTheme(matched[0].Name, matched[0].Color), nil
}

func PaletteThemes(pal string) (Themes, error) {
	p, err := getPalette(pal)
	if err != nil {
		return nil, err
	}
	colors := paletteColors(*p)
	return lo.Map(colors, func(c gamut.Color, _ int) *Theme {
		return ColorTheme(c.Name, c.Color)
	}), nil
}

func PaletteRandomThemes(pal string, num int) (Themes, error) {
	p, err := getPalette(pal)
	if err != nil {
		return nil, err
	}
	pc := p.Colors()
	var ret Themes = lo.Times(num, func(_ int) *Theme {
		return randomTheme(pc)
	})
	return ret, nil
}

func getPalette(pal string) (*gamut.Palette, error) {
	switch pal {
	case "crayola", "":
		return &palette.Crayola, nil
	case "css":
		return &palette.CSS, nil
	case "ral":
		return &palette.RAL, nil
	case "resene":
		return &palette.Resene, nil
	case "wikipedia":
		return &palette.Wikipedia, nil
	default:
		return nil, errors.Errorf("no palette available with key [%s]", pal)
	}
}

func paletteColors(p gamut.Palette) gamut.Colors {
	ret := p.Colors()
	slices.SortFunc(ret, func(l gamut.Color, r gamut.Color) int {
		return cmp.Compare(l.Name, r.Name)
	})
	return ret
}

func randomTheme(colors gamut.Colors) *Theme {
	x := colors[util.RandomInt(len(colors))]
	return ColorTheme(x.Name, x.Color)
}
