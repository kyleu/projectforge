// Content managed by Project Forge, see [projectforge.md] for details.
package theme

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Theme struct {
	Key   string  `json:"-"`
	Light *Colors `json:"light"`
	Dark  *Colors `json:"dark"`
	css   string
}

func (t *Theme) CSS(indent int) string {
	if t.css != "" {
		return t.css
	}
	sb := &strings.Builder{}
	sb.WriteString("/* theme: " + t.Key + " */\n")
	sb.WriteString(t.Light.CSS(":root", indent))
	sb.WriteString(t.Light.CSS(".mode-light", indent))
	sb.WriteString(t.Dark.CSS(".mode-dark", indent))
	addLine(sb, "", indent)
	addLine(sb, "@media (prefers-color-scheme: dark) {", indent)
	sb.WriteString(t.Dark.CSS(":root", indent+1))
	sb.WriteString(t.Light.CSS(".mode-light", indent+1))
	sb.WriteString(t.Dark.CSS(".mode-dark", indent+1))
	addLine(sb, "}", indent)
	t.css = sb.String()
	return t.css
}

func (t *Theme) Clone(key string) *Theme {
	return &Theme{Key: key, Light: t.Light.Clone(), Dark: t.Dark.Clone()}
}

func (t *Theme) Equals(x *Theme) bool {
	return t.Light.Equals(x.Light) && t.Dark.Equals(x.Dark)
}

func (t *Theme) ToGo() string {
	var ret []string
	add := func(ind int, s string, args ...any) {
		ret = append(ret, util.StringRepeat("\t", ind+1)+fmt.Sprintf(s, args...))
	}
	addColors := func(c *Colors) {
		add(2, "Border: %q, LinkDecoration: %q,", c.Border, c.LinkDecoration)
		add(2, "Foreground: %q, ForegroundMuted: %q,", c.Foreground, c.ForegroundMuted)
		add(2, "Background: %q, BackgroundMuted: %q,", c.Background, c.BackgroundMuted)
		add(2, "LinkForeground: %q, LinkVisitedForeground: %q,", c.LinkForeground, c.LinkVisitedForeground)
		add(2, "NavForeground: %q, NavBackground: %q,", c.NavForeground, c.NavBackground)
		add(2, "MenuForeground: %q, MenuSelectedForeground: %q,", c.MenuForeground, c.MenuSelectedForeground)
		add(2, "MenuBackground: %q, MenuSelectedBackground: %q,", c.MenuBackground, c.MenuSelectedBackground)
		add(2, "ModalBackdrop: %q, Success: %q, Error: %q,", c.ModalBackdrop, c.Success, c.Error)
	}
	add(0, "&Theme{")
	add(1, "Key: %q,", t.Key)
	add(1, "Light: &Colors{")
	addColors(t.Light)
	add(1, "},")
	add(1, "Dark: &Colors{")
	addColors(t.Dark)
	add(1, "},")
	add(0, "},")
	return strings.Join(ret, util.StringDefaultLinebreak)
}

type Themes []*Theme

func (t Themes) Sort() Themes {
	slices.SortFunc(t, func(l *Theme, r *Theme) int {
		if l.Key == ThemeDefault.Key {
			return 0
		}
		if r.Key == ThemeDefault.Key {
			return 0
		}
		return cmp.Compare(strings.ToLower(l.Key), strings.ToLower(r.Key))
	})
	return t
}

func (t Themes) Replace(n *Theme) Themes {
	for idx, o := range t {
		if o.Key == n.Key {
			t[idx] = n
			return t
		}
	}
	ret := append(Themes{}, t...)
	ret = append(ret, n)
	return ret.Sort()
}

func (t Themes) Contains(key string) bool {
	return lo.ContainsBy(t, func(x *Theme) bool {
		return x.Key == key
	})
}

func (t Themes) Get(key string) *Theme {
	return lo.FindOrElse(t, nil, func(x *Theme) bool {
		return x.Key == key
	})
}

func (t Themes) Remove(key string) Themes {
	return lo.Filter(t, func(thm *Theme, _ int) bool {
		return thm.Key != key
	})
}

func addLine(sb io.StringWriter, s string, indent int) {
	indention := util.StringRepeat("  ", indent)
	_, _ = sb.WriteString(indention + s + "\n")
}
