package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"$PF_PACKAGE$/app/util"
)

func rootF(*cobra.Command, []string) error {
	return startServer(_flags)
}

func rootCmd() *cobra.Command {
	ret := &cobra.Command{Use: util.AppKey, Short: util.AppSummary, RunE: rootF}

	// $PF_SECTION_START(cmds)$
	ret.AddCommand(serverCmd(), siteCmd(), allCmd())
	// $PF_SECTION_END(cmds)$

	ret.PersistentFlags().StringVarP(&_flags.ConfigDir, "dir", "d", "", "directory for configuration, defaults to system config dir")
	ret.PersistentFlags().BoolVarP(&_flags.Debug, "verbose", "v", false, "enables verbose logging and additional checks")
	ret.PersistentFlags().StringVarP(&_flags.Address, "addr", "a", "127.0.0.1", "address to listen on, defaults to [127.0.0.1]")
	ret.PersistentFlags().Uint16VarP(&_flags.Port, "port", "p", util.AppPort, fmt.Sprintf("port to listen on, defaults to [%d]", util.AppPort))

	return ret
}
