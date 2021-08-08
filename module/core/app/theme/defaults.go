package theme

var ThemeDefault = &Theme{
	Key: "default",
	Light: &Colors{
		Border: "{{{ .Theme.Light.Border }}}", LinkDecoration: "{{{ .Theme.Light.LinkDecoration }}}",
		Foreground: "{{{ .Theme.Light.Foreground }}}", ForegroundMuted: "{{{ .Theme.Light.ForegroundMuted }}}",
		Background: "{{{ .Theme.Light.Background }}}", BackgroundMuted: "{{{ .Theme.Light.BackgroundMuted }}}",
		LinkForeground: "{{{ .Theme.Light.LinkForeground }}}", LinkVisitedForeground: "{{{ .Theme.Light.LinkVisitedForeground }}}",
		NavForeground: "{{{ .Theme.Light.NavForeground }}}", NavBackground: "{{{ .Theme.Light.NavBackground }}}",
		MenuForeground: "{{{ .Theme.Light.MenuForeground }}}", MenuBackground: "{{{ .Theme.Light.MenuBackground }}}", MenuSelectedBackground: "{{{ .Theme.Light.MenuSelectedBackground }}}", MenuSelectedForeground: "{{{ .Theme.Light.MenuSelectedForeground }}}",
		ModalBackdrop: "{{{ .Theme.Light.ModalBackdrop }}}", Success: "{{{ .Theme.Light.Success }}}", Error: "{{{ .Theme.Light.Error }}}",
	},
	Dark: &Colors{
		Border: "{{{ .Theme.Dark.Border }}}", LinkDecoration: "{{{ .Theme.Dark.LinkDecoration }}}",
		Foreground: "{{{ .Theme.Dark.Foreground }}}", ForegroundMuted: "{{{ .Theme.Dark.ForegroundMuted }}}",
		Background: "{{{ .Theme.Dark.Background }}}", BackgroundMuted: "{{{ .Theme.Dark.BackgroundMuted }}}",
		LinkForeground: "{{{ .Theme.Dark.LinkForeground }}}", LinkVisitedForeground: "{{{ .Theme.Dark.LinkVisitedForeground }}}",
		NavForeground: "{{{ .Theme.Dark.NavForeground }}}", NavBackground: "{{{ .Theme.Dark.NavBackground }}}",
		MenuForeground: "{{{ .Theme.Dark.MenuForeground }}}", MenuBackground: "{{{ .Theme.Dark.MenuBackground }}}", MenuSelectedBackground: "{{{ .Theme.Dark.MenuSelectedBackground }}}", MenuSelectedForeground: "{{{ .Theme.Dark.MenuSelectedForeground }}}",
		ModalBackdrop: "{{{ .Theme.Dark.ModalBackdrop }}}", Success: "{{{ .Theme.Dark.Success }}}", Error: "{{{ .Theme.Dark.Error }}}",
	},
}
