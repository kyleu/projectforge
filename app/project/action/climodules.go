package action

import (
	"fmt"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func promptModules(current []string, mods module.Modules) []string {
	modKeys := util.ArraySorted(mods.Keys())
	const msg = "Enter the modules your project will use as a comma-separated list (or \"all\") from choices"
	modPromptString := fmt.Sprintf("%s:\n  %s", msg, util.StringJoin(modKeys, ", "))
	modStrings := promptString(modPromptString, util.StringJoin(current, ", "))
	ret := util.StringSplitAndTrim(modStrings, ",")
	if (len(ret) == 1 && ret[0] == "core") && ret[0] == "all" {
		return modKeys
	}
	return ret
}
