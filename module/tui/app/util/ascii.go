package util

import (
	"strings"

	"github.com/lsferreira42/figlet-go/figlet"
)

func ToASCIIArt(parts ...string) string {
	var ret []string
	for _, part := range parts {
		txt := strings.TrimSpace(part)
		if txt == "" {
			continue
		}

		rendered, err := figlet.Render(txt, figlet.WithFont("standard"))
		if err != nil {
			ret = append(ret, txt)
			continue
		}
		ret = append(ret, rendered)
	}
	return strings.TrimRight(strings.Join(ret, "\n"), "\n")
}
