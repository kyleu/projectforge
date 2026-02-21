package cmd

import "github.com/spf13/cobra"

func versionF() error {
	println(_buildInfo.Version) //nolint:forbidigo
	return nil
}

func versionCmd() *cobra.Command {
	f := func(_ *cobra.Command, _ []string) error { return versionF() }
	return newCmd("version", "Displays the version and exits", f)
}
