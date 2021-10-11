package checks

import (
	"github.com/kyleu/projectforge/app/doctor"
	"github.com/kyleu/projectforge/app/util"
)

var AllChecks = doctor.Checks{air, golang, imagemagick, make, node, qtc}

func GetCheck(key string) *doctor.Check {
	for _, x := range AllChecks {
		if x.Key == key {
			return x
		}
	}
	return nil
}

func ForModules(modules []string) doctor.Checks {
	var ret doctor.Checks
	for _, c := range AllChecks {
		hit := len(c.Modules) == 0
		for _, mod := range c.Modules {
			if util.StringArrayContains(modules, mod) {
				hit = true
				break
			}
		}
		if !hit {
			continue
		}
		ret = append(ret, c)
	}
	return ret
}

func CheckAll(modules []string) doctor.Results {
	var ret doctor.Results
	for _, c := range ForModules(modules) {
		ret = append(ret, c.Check())
	}
	return ret
}
