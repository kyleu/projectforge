package cmd

import (
	"context"

	"github.com/kyleu/projectforge/app/action"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func actionF(ctx context.Context, t action.Type, args []string) error {
	if err := initIfNeeded(); err != nil {
		return errors.Wrap(err, "error initializing application")
	}

	_, cfg := extractConfig(args)
	projectKey := "TODO"

	return runToCompletion(ctx, projectKey, t, cfg)
}

func actionCmd(ctx context.Context, t action.Type) *cobra.Command {
	f := func(cmd *cobra.Command, args []string) error { return actionF(ctx, t, args) }
	return &cobra.Command{Use: t.Key, Short: t.Description, RunE: f}
}

func actionCommands() []*cobra.Command {
	ctx := context.Background()
	ret := make([]*cobra.Command, 0, len(action.AllTypes))
	for _, a := range action.AllTypes {
		if !a.Hidden {
			ret = append(ret, actionCmd(ctx, a))
		}
	}
	return ret
}
