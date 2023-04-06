package checks

import (
	"context"

	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

var AllChecks = doctor.Checks{pf, prj, repo, air, git, golang, imagemagick, mke, node, qtc}

func GetCheck(key string) *doctor.Check {
	return AllChecks.Get(key)
}

func ForModules(modules []string) doctor.Checks {
	var ret doctor.Checks
	for _, c := range AllChecks {
		hit := len(c.Modules) == 0
		for _, mod := range c.Modules {
			if slices.Contains(modules, mod) {
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

func CheckAll(ctx context.Context, modules []string, logger util.Logger) doctor.Results {
	ctx, span, logger := telemetry.StartSpan(ctx, "doctor:checkall", logger)
	defer span.Complete()

	var ret doctor.Results
	for _, c := range ForModules(modules) {
		ret = append(ret, c.Check(ctx, logger))
	}
	return ret
}

func noop(_ context.Context, r *doctor.Result, _ string) *doctor.Result {
	return r
}
