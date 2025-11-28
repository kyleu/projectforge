package action

import (
	"path"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type CustomCmdResult struct {
	Cmd      string   `json:"cmd,omitzero"`
	Exit     int      `json:"exit,omitzero"`
	Output   []string `json:"output,omitzero"`
	Duration int      `json:"duration,omitzero"`
}

func runCustom(p *project.Project, cmd string, dir string, logger util.Logger) (*CustomCmdResult, error) {
	if cmd == "" {
		return nil, errors.New("no command provided")
	}
	pth := p.Path
	if dir != "" {
		pth = path.Join(p.Path, dir)
	}
	logger.Debugf("running custom command [%s] in [%s]...", cmd, pth)
	t := util.TimerStart()
	exit, x, err := util.RunProcessSimple(cmd, pth)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to run custom command [%s] in [%s]", cmd, pth)
	}
	return &CustomCmdResult{Cmd: cmd, Exit: exit, Output: util.StringSplitLines(x), Duration: t.End()}, nil
}
