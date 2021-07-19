package cmd

import (
	"$PF_PACKAGE$/app"
	"go.uber.org/zap"
)

func Run(bi *app.BuildInfo) (*zap.SugaredLogger, error) {
	_buildInfo = bi

	c := rootCmd()
	// $PF_SECTION_START(cmds)$
	c.AddCommand(/* add new commands here */)
	// $PF_SECTION_END(cmds)$
	if err := c.Execute(); err != nil {
		return _logger, err
	}
	return _logger, nil
}
