package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"{{{ .Package }}}/app/util"
)

var rootCtx = context.Background()

func rootF(*cobra.Command, []string) error {
	// $PF_SECTION_START(rootAction)$
	return startServer(rootCtx, _flags)
	// $PF_SECTION_END(rootAction)$
}

func rootCmd(ctx context.Context) *cobra.Command {
	short := fmt.Sprintf("%s %s - %s", util.AppName, _buildInfo.Version, util.AppSummary)
	ret := newCmd(util.AppKey, short, rootF)
	ret.AddCommand(serverCmd(){{{ if .HasModule "marketing" }}}, siteCmd(), allCmd(){{{ end }}}{{{ if .HasModule "mcp" }}}, mcpCmd(){{{ end }}}{{{ if .HasModule "migration" }}}, migrateCmd(){{{ end }}}{{{ if .HasModule "upgrade" }}}, upgradeCmd(){{{ end }}}{{{ if .HasModule "wasmserver" }}}, wasmCmd(){{{ end }}})
	// $PF_SECTION_START(cmds)$ - Add your commands here by calling ret.AddCommand
	// $PF_SECTION_END(cmds)$
	ret.AddCommand(versionCmd())

	ret.PersistentFlags().StringVarP(&_flags.WorkingDir, "working_dir", "w", ".", "directory for projects, defaults to current dir")
	ret.PersistentFlags().StringVarP(&_flags.ConfigDir, "config_dir", "c", "", "directory for configuration, defaults to system config dir")
	ret.PersistentFlags().BoolVarP(&_flags.Debug, "verbose", "v", false, "enables verbose logging and additional checks")
	ret.PersistentFlags().StringVarP(&_flags.Address, "addr", "a", "127.0.0.1", "address to listen on, defaults to [127.0.0.1]")
	ret.PersistentFlags().Uint16VarP(&_flags.Port, "port", "p", util.AppPort, fmt.Sprintf("port to listen on, defaults to [%d]", util.AppPort))

	return ret
}
