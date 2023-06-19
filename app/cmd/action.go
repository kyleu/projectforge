package cmd

import (
	"context"
	"strings"

	"github.com/muesli/coral"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"

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
	if t.Key == action.TypeCreate.Key {
		aliases = []string{"init", "new"}
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
		for _, e := range r.Errors {
			_logger.Warn(" - " + e)
		}
	}
	if r.Modules.DiffCount(false) > 0 {
		for _, m := range r.Modules {
			for _, d := range m.DiffsFiltered(false) {
				_logger.Infof("%s [%s]:", d.Path, d.Status)
				for _, c := range d.Changes {
					for _, l := range c.Lines {
						_logger.Info(strings.TrimSuffix(l.String(), "\n"))
					}
				}
			}
		}
	}
	if r.Data != nil {
		deps, ok := r.Data.(build.Dependencies)
		if ok {
			for _, dep := range deps {
				_logger.Infof(dep.String())
			}
		}
	}
}
