// Package cmd - Content managed by Project Forge, see [projectforge.md] for details.
package cmd

import "github.com/muesli/coral"

func versionF() error {
	println(_buildInfo.Version) //nolint:forbidigo
	return nil
}

func versionCmd() *coral.Command {
	f := func(_ *coral.Command, _ []string) error { return versionF() }
	return &coral.Command{Use: "version", Short: "Displays the version and exits", RunE: f}
}
