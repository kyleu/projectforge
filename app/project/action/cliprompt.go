//go:build !js
// +build !js

package action

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/util"
)

func promptString(label string, query string, curr string) (string, error) {
	promptTotal++
	title := fmt.Sprintf("%d: %s", promptTotal, query)
	if promptTotal > 1 {
		title = util.StringDefaultLinebreak + title
	}
	text := curr
	in := huh.NewInput().Title(title).Value(&text)
	if curr == "" {
		in = in.Description("Optional")
	} else {
		in = in.Description("Default: " + curr)
	}
	f := huh.NewForm(huh.NewGroup(in))
	err := f.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", errors.New("project creation canceled")
		}
		clilog("error: " + err.Error() + util.StringDefaultLinebreak)
		return "", err
	}
	ret := util.OrDefault(strings.TrimSpace(text), curr)
	logPromptAnswer(label, ret)
	return ret, nil
}
