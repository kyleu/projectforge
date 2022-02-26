package theme

import (
	"sort"

	"github.com/kyleu/projectforge/app/util"
	"github.com/muesli/gamut"
	"github.com/muesli/gamut/palette"
	"github.com/pkg/errors"
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
	ret := make(Themes, 0, len(colors))
	for _, c := range colors {
		ret = append(ret, ColorTheme(c.Name, c.Color))
	}
	return ret, nil
}

func PaletteRandomThemes(pal string, num int) (Themes, error) {
	ret := make(Themes, 0, num)
	p, err := getPalette(pal)
	if err != nil {
		return nil, err
	}
	pc := p.Colors()
	for i := 0; i < num; i++ {
		ret = append(ret, randomTheme(pc))
	}
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
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Name < ret[j].Name
	})
	return ret
}

func randomTheme(colors gamut.Colors) *Theme {
	x := colors[util.RandomInt(len(colors))]
	return ColorTheme(x.Name, x.Color)
}
