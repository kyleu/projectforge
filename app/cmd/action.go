package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/project/build"
	"projectforge.dev/projectforge/app/util"
)

func actionF(ctx context.Context, t action.Type, args []string) error {
	logger, err := initIfNeeded(ctx)
	if err != nil {
		return errors.Wrap(err, "error initializing application")
	}
	cfg := extractConfig(args)
	logResult(t, runToCompletion(ctx, "", t, cfg, logger), logger)
	return nil
}

func actionCmd(ctx context.Context, t action.Type) *cobra.Command {
	f := func(_ *cobra.Command, args []string) error { return actionF(ctx, t, args) }
	var aliases []string
	switch t.Key {
	case action.TypeCreate.Key:
		aliases = []string{"init", "new"}
	case action.TypeBuild.Key:
		aliases = []string{"make"}
	case action.TypeGenerate.Key:
		aliases = []string{"gen"}
	}
	ret := newCmd(t.Key, t.Description, f, aliases...)
	if t.Matches(action.TypeBuild) {
		lo.ForEach(action.AllBuilds, func(x *action.Build, _ int) {
			k := x.Key
			fx := func(_ *cobra.Command, args []string) error {
				return actionF(ctx, t, append(util.ArrayCopy(args), k))
			}
			ret.AddCommand(newCmd(x.Key, x.Description, fx))
		})
	}
	return ret
}

func actionCommands(ctx context.Context) []*cobra.Command {
	return lo.FilterMap(action.AllTypes, func(a action.Type, _ int) (*cobra.Command, bool) {
		if a.Hidden {
			return nil, false
		}
		return actionCmd(ctx, a), true
	})
}

func logResult(t action.Type, r *action.Result, logger util.Logger) {
	logger.Infof("%s [%s]: %s in [%s]", util.AppName, t.String(), r.Status, util.MicrosToMillis(r.Duration))
	if len(r.Errors) > 0 {
		logger.Warnf("Errors:")
		lo.ForEach(r.Errors, func(e string, _ int) {
			logger.Warn(" - " + e)
		})
	}
	if r.Modules.DiffCount(false) > 0 {
		lo.ForEach(r.Modules, func(m *module.Result, _ int) {
			lo.ForEach(m.DiffsFiltered(false), func(d *diff.Diff, _ int) {
				logger.Infof("%s [%s]:", d.Path, d.Status)
				lo.ForEach(d.Changes, func(c *diff.Change, _ int) {
					lo.ForEach(c.Lines, func(l *diff.Line, _ int) {
						rs := util.Str(l.String())
						logger.Info(rs.TrimSuffix("\n").TrimSuffix("\r"))
					})
				})
			})
		})
	}
	if r.Data != nil {
		deps, err := util.Cast[build.Dependencies](r.Data)
		if err == nil {
			lo.ForEach(deps, func(dep *build.Dependency, _ int) {
				logger.Infof(dep.String())
			})
		}
	}
}
