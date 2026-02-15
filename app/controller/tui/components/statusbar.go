package components

import (
	"strings"

	"projectforge.dev/projectforge/app/controller/tui/style"
)

func RenderStatus(status string, err string, shortHelp []string, st style.Styles) string {
	var top string
	if err != "" {
		top = st.Error.Render(err)
	} else {
		top = st.Status.Render(status)
	}
	help := st.Muted.Render(strings.Join(shortHelp, " | "))
	return top + "\n" + help
}
