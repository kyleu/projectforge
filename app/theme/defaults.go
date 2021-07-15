package theme

var ThemeDefault = &Theme{
	Key: "default",
	Light: &Colors{
		Border: "1px solid #dddddd", LinkDecoration: "none",
		Foreground: "#000000", ForegroundMuted: "#999999",
		Background: "#ffffff", BackgroundMuted: "#eeeeee",
		LinkForeground: "#2d414e", LinkVisitedForeground: "#406379",
		NavForeground: "#000000", NavBackground: "#4f9abd",
		MenuForeground: "#000000", MenuBackground: "#f0f8ff", MenuSelectedBackground: "#faebd7", MenuSelectedForeground: "#000000",
		ModalBackdrop: "rgba(77, 77, 77, .7)", Success: "#008000", Error: "#FF0000",
	},
	Dark: &Colors{
		Border: "1px solid #666666", LinkDecoration: "none",
		Foreground: "#ffffff", ForegroundMuted: "#999999",
		Background: "#121212", BackgroundMuted: "#333333",
		LinkForeground: "#dddddd", LinkVisitedForeground: "#aaaaaa",
		NavForeground: "#ffffff", NavBackground: "#2d414e",
		MenuForeground: "#dddddd", MenuBackground: "#171f24", MenuSelectedBackground: "#333333", MenuSelectedForeground: "#ffffff",
		ModalBackdrop: "rgba(33, 33, 33, .7)", Success: "#008000", Error: "#FF0000",
	},
}
