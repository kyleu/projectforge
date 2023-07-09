package checks

import (
	"context"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/lib/telemetry"
	"projectforge.dev/projectforge/app/util"
)

var AllChecks = doctor.Checks{
	pf, homebrew, choco, golang, mke, node, git, air, qtc, imagemagick, inkscape, repo, prj,
}

func GetCheck(key string) *doctor.Check {
	return AllChecks.Get(key)
}

func ForModules(modules []string) doctor.Checks {
	var ret doctor.Checks
	lo.ForEach(AllChecks, func(c *doctor.Check, _ int) {
		if !c.Applies() {
			return
		}
		hit := len(c.Modules) == 0 || lo.ContainsBy(c.Modules, func(mod string) bool {
			return lo.Contains(modules, mod)
		})
		if !hit {
			return
		}
		ret = append(ret, c)
	})
	return ret
}

func CheckAll(ctx context.Context, modules []string, logger util.Logger, exclude ...string) doctor.Results {
	ctx, span, logger := telemetry.StartSpan(ctx, "doctor:checkall", logger)
	defer span.Complete()

	return lo.FilterMap(ForModules(modules), func(c *doctor.Check, _ int) (*doctor.Result, bool) {
		if lo.Contains(exclude, c.Key) {
			return nil, false
		}
		return c.Check(ctx, logger), true
	})
}

func noop(_ context.Context, r *doctor.Result, _ string) *doctor.Result {
	return r
}
