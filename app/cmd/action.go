package cmd

import (
	"context"
	"strings"

	"github.com/muesli/coral"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/file/diff"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project/action"
	"projectforge.dev/projectforge/app/project/build"
	"projectforge.dev/projectforge/app/util"
)

func actionF(ctx context.Context, t action.Type, args []string) error {
	if err := initIfNeeded(); err != nil {
		return errors.Wrap(err, "error initializing application")
	}
	cfg := extractConfig(args)
	logResult(t, runToCompletion(ctx, "", t, cfg))
	return nil
}

func actionCmd(ctx context.Context, t action.Type) *coral.Command {
	f := func(cmd *coral.Command, args []string) error { return actionF(ctx, t, args) }
	var aliases []string
	switch t.Key {
	case action.TypeCreate.Key:
		aliases = []string{"init", "new"}
	case action.TypeBuild.Key:
		aliases = []string{"make"}
	case action.TypeGenerate.Key:
		aliases = []string{"gen"}
	}
	ret := &coral.Command{Use: t.Key, Short: t.Description, RunE: f, Aliases: aliases}
	if t.Key == action.TypeBuild.Key {
		lo.ForEach(action.AllBuilds, func(x *action.Build, _ int) {
			k := x.Key
			fx := func(cmd *coral.Command, args []string) error {
				return actionF(ctx, t, append(slices.Clone(args), k))
			}
			ret.AddCommand(&coral.Command{Use: x.Key, Short: x.Description, RunE: fx})
		})
	}
	return ret
}

func actionCommands() []*coral.Command {
	ctx := context.Background()
	ret := lo.FilterMap(action.AllTypes, func(a action.Type, _ int) (*coral.Command, bool) {
		if a.Hidden {
			return nil, false
		}
		return actionCmd(ctx, a), true
	})
	return ret
}

func logResult(t action.Type, r *action.Result) {
	_logger.Infof("%s [%s]: %s in [%s]", util.AppName, t.String(), r.Status, util.MicrosToMillis(r.Duration))
	if len(r.Errors) > 0 {
		_logger.Warnf("Errors:")
		lo.ForEach(r.Errors, func(e string, _ int) {
			_logger.Warn(" - " + e)
		})
	}
	if r.Modules.DiffCount(false) > 0 {
		lo.ForEach(r.Modules, func(m *module.Result, _ int) {
			lo.ForEach(m.DiffsFiltered(false), func(d *diff.Diff, _ int) {
				_logger.Infof("%s [%s]:", d.Path, d.Status)
				lo.ForEach(d.Changes, func(c *diff.Change, _ int) {
					lo.ForEach(c.Lines, func(l *diff.Line, _ int) {
						_logger.Info(strings.TrimSuffix(strings.TrimSuffix(l.String(), "\n"), "\r"))
					})
				})
			})
		})
	}
	if r.Data != nil {
		deps, ok := r.Data.(build.Dependencies)
		if ok {
			lo.ForEach(deps, func(dep *build.Dependency, _ int) {
				_logger.Infof(dep.String())
			})
		}
	}
}
