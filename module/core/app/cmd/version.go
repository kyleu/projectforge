package cmd

import "github.com/muesli/coral"

func versionF() error {
	println(_buildInfo.Version) //nolint:forbidigo
	return nil
}

func versionCmd() *coral.Command {
	f := func(_ *coral.Command, _ []string) error { return versionF() }
	return newCmd("version", "Displays the version and exits", f)
}
