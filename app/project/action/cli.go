package action

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

var promptTotal = 1

var (
	cliInfoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	cliSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	cliErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	cliLabelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	cliValueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
)

func clilog(level string, s string) {
	style := cliInfoStyle
	switch level {
	case "error":
		style = cliErrorStyle
	case "success":
		style = cliSuccessStyle
	}
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if line != "" {
			_, _ = os.Stdout.WriteString(style.Render(line))
		}
		_, _ = os.Stdout.WriteString("\n")
	}
}

const promptLabelWidth = len("Author GitHub")

func logPromptAnswer(label string, value string) {
	pad := promptLabelWidth - len(label)
	if pad < 0 {
		pad = 0
	}
	spacer := strings.Repeat(" ", pad+1)
	clilog("info", cliLabelStyle.Render(label+":")+spacer+cliValueStyle.Render(value))
}

func cliProject(ctx context.Context, p *project.Project, mods module.Modules, logger util.Logger) error {
	clilog("success", util.AppName)
	clilog("info", "Let's create a new project!\n")
	clilog("debug", "Checking a few things...")
	errResults := checks.CheckAll(ctx, p.Modules, logger, checks.Core(false).Keys()...).Errors()
	if len(errResults) > 0 {
		msgs := lo.Map(errResults, func(r *doctor.Result, _ int) string {
			return r.ErrorString()
		})
		return errors.New(util.StringJoin(msgs, util.StringDefaultLinebreak))
	}
	targetPath := p.Path
	if absPath, err := filepath.Abs(targetPath); err == nil {
		targetPath = absPath
	}
	clilog("success", "All good, project will be created in ["+targetPath+"]\n")

	if err := cliGather(ctx, p, mods, logger); err != nil {
		return err
	}
	return nil
}
