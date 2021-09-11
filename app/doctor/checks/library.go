package checks

import (
	"github.com/kyleu/projectforge/app/doctor"
	"github.com/kyleu/projectforge/app/util"
)

var AllChecks = doctor.Checks{node, rsvg}

func CheckAll(modules []string) doctor.Results {
	var ret doctor.Results
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
		r := doctor.NewResult(c, c.Key, c.Title, c.Summary)
		ret = append(ret, c.Fn(r))
	}
	return ret
}
