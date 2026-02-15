//go:build !js && !aix
// +build !js,!aix

package action

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func promptModules(current []string, mods module.Modules) ([]string, error) {
	promptTotal++
	title := fmt.Sprintf("%d: Select the modules your project will use", promptTotal)
	if promptTotal > 1 {
		title = util.StringDefaultLinebreak + title
	}
	ret := append([]string{}, current...)
	var opts []huh.Option[string]
	for _, mod := range mods.Sorted() {
		opts = append(opts, huh.NewOption(fmt.Sprintf("%s: %s", mod.Key, mod.Description), mod.Key))
	}
	f := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title(title).
			Description("Use arrow keys to move and space to toggle modules").
			Options(opts...).
			Value(&ret),
	)).WithProgramOptions(tea.WithAltScreen())
	if err := f.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, errors.New("project creation canceled")
		}
		clilog("error: " + err.Error() + util.StringDefaultLinebreak)
		return nil, err
	}
	ret = util.ArraySorted(util.ArrayRemoveDuplicates(ret))
	logPromptAnswer("Modules", strings.Join(ret, ", "))
	return ret, nil
}
