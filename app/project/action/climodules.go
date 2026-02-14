package action

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func promptModules(current []string, mods module.Modules) ([]string, error) {
	promptTotal++
	ret := append([]string{}, current...)
	var opts []huh.Option[string]
	for _, mod := range mods.Sorted() {
		opts = append(opts, huh.NewOption(fmt.Sprintf("%s: %s", mod.Key, mod.Description), mod.Key))
	}
	f := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title(fmt.Sprintf("%d: Select the modules your project will use", promptTotal)).
			Description("Use arrow keys to move and space to toggle modules").
			Options(opts...).
			Value(&ret),
	))
	if err := f.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, errors.New("project creation canceled")
		}
		clilog("error: " + err.Error() + util.StringDefaultLinebreak)
		return nil, err
	}
	ret = util.ArraySorted(util.ArrayRemoveDuplicates(ret))
	return ret, nil
}
