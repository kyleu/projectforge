package theme

import (
	"{{{ .Package }}}/app/util"
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
