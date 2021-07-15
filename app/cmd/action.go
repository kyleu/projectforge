package cmd

import (
	"github.com/kyleu/projectforge/app/action"
	"github.com/spf13/cobra"
)

func actionF(t action.Type, args []string) error {
	_, cfg := extractConfig(args)
	return runToCompletion(t, cfg)
}

func actionCmd(t action.Type) *cobra.Command {
	f := func(cmd *cobra.Command, args []string) error { return actionF(t, args) }
	return &cobra.Command{Use: t.Key, Short: t.Description, RunE: f}
}

func actionCommands() []*cobra.Command {
	ret := make([]*cobra.Command, 0, len(action.AllTypes))
	for _, a := range action.AllTypes {
		ret = append(ret, actionCmd(a))
	}
	return ret
}
