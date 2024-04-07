package cmd

import (
	"github.com/muesli/coral"

	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/tools/pfdb/app/export"
)

func rootF(*coral.Command, []string) error {
	svc, err := export.NewService(_logger)
	if err != nil {
		return err
	}
	return svc.Run(_flags.DatabaseType, _flags.DatabaseKey)
}

func rootCmd() *coral.Command {
	short := "Project Forge Database Tools"
	ret := &coral.Command{Use: util.AppKey, Short: short, RunE: rootF}

	ret.PersistentFlags().BoolVarP(&_flags.Debug, "verbose", "v", false, "enables verbose logging and additional checks")
	ret.PersistentFlags().StringVarP(&_flags.DatabaseType, "dbtype", "d", "postgres", "database type to use")
	ret.PersistentFlags().StringVarP(&_flags.DatabaseKey, "dbkey", "k", "postgres", "database key to load configuration from")

	return ret
}
