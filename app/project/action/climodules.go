package action

import (
	"fmt"
	"slices"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/util"
)

func promptModules(current []string, mods module.Modules) []string {
	modKeys := util.ArraySorted(mods.Keys())
	const msg = "Enter the modules your project will use as a comma-separated list\n" +
		"  Type \"list\" to see the available choices, or \"all\" for everything."
	modStrings := promptString(msg, util.StringJoin(current, ", "))
	ret := util.StringSplitAndTrim(modStrings, ",")
	if len(ret) == 1 {
		switch ret[0] {
		case "all":
			return modKeys
		case "list":
			clilog("Available modules:")
			for _, mod := range mods.Sorted() {
				clilog(fmt.Sprintf(util.StringDefaultLinebreak+"  %12s %s", mod.Key+":", mod.Description))
			}
			clilog(util.StringDefaultLinebreak + util.StringDefaultLinebreak)
			promptTotal--
			return promptModules(current, mods)
		}
	}
	ret = util.ArraySorted(util.ArrayRemoveDuplicates(ret))
	for _, mod := range ret {
		if !slices.Contains(modKeys, mod) {
			clilog("error: unknown module: " + mod + util.StringDefaultLinebreak)
			return promptModules(current, mods)
		}
	}
	return ret
}
