package checks

import (
	"projectforge.dev/app/doctor"
	"projectforge.dev/app/util"
	"go.uber.org/zap"
)

var AllChecks = doctor.Checks{pf, prj, repo, air, git, golang, imagemagick, mke, node, qtc}

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

func CheckAll(modules []string, logger *zap.SugaredLogger) doctor.Results {
	var ret doctor.Results
	for _, c := range ForModules(modules) {
		ret = append(ret, c.Check(logger))
	}
	return ret
}

func noop(r *doctor.Result, out string) *doctor.Result {
	return r
}
