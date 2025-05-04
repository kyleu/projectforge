package action

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func clilog(s string) {
	print(s) //nolint:forbidigo
}

func cliProject(ctx context.Context, p *project.Project, modKeys []string, logger util.Logger) error {
	clilog(util.AppName + "\nLet's create a new project!\n\n")

	clilog("Checking a few things...\n")
	errResults := checks.CheckAll(ctx, p.Modules, logger, checks.Core(false).Keys()...).Errors()
	if len(errResults) > 0 {
		msgs := lo.Map(errResults, func(r *doctor.Result, _ int) string {
			return r.ErrorString()
		})
		return errors.New(util.StringJoin(msgs, util.StringDefaultLinebreak))
	}
	clilog("All good, project will be created in the current directory\n\n")
	if p.Key == "" {
		path, _ := os.Getwd()
		_, path = util.StringSplitPath(path)
		p.Key = path
	}
	p.Key = promptString("Enter a project key; must only contain alphanumerics", strings.ToLower(p.Key))

	if p.Name == "" {
		p.Name = p.Key
	}
	p.Name = promptString("Enter a project name; use title case and spaces if needed", p.Name)

	if p.Icon == "" {
		p.Icon = "star"
	}

	if p.Exec == "" {
		p.Exec = p.Key
	}

	p.Version = promptString("Enter a version, such as 0.0.0", p.Version)

	gatherProjectInfo(p)

	if p.Info.Summary == "" {
		p.Info.Summary = "A simple project"
	}
	p.Info.Summary = promptString("Enter a one-line description of this project", p.Info.Summary)

	p.Info.AuthorName = promptString("Enter the name of this project's owner", p.Info.AuthorName)
	if p.Info.AuthorEmail == "" {
		p.Info.AuthorEmail = fmt.Sprintf("dev@%s.com", p.Key)
	}
	p.Info.AuthorEmail = promptString("Enter the email address of this project's owner", p.Info.AuthorEmail)
	p.Info.AuthorID = promptString("Enter the GitHub username(s) of this project's owner", p.Info.AuthorID)
	p.Info.Team = promptString("Enter the team that owns this project", p.Info.Team)

	if p.Port == 0 {
		p.Port = 20000
	}
	prt, _ := strconv.ParseInt(promptString("Enter the default port your http server will run on", fmt.Sprint(p.Port)), 10, 32)
	p.Port = int(prt)

	const msg = "Enter the modules your project will use as a comma-separated list (or \"all\") from choices"
	modPromptString := fmt.Sprintf("%s:\n  %s", msg, util.StringJoin(modKeys, ", "))
	mods := promptString(modPromptString, util.StringJoin(p.Modules, ", "))
	p.Modules = util.StringSplitAndTrim(mods, ",")
	if len(p.Modules) == 1 && p.Modules[0] == "all" {
		p.Modules = modKeys
	}

	if p.Info.License == "" {
		p.Info.License = "Proprietary"
	}
	p.Info.License = promptString("Enter the license used by this project", p.Info.License)

	return nil
}

func gatherProjectInfo(p *project.Project) {
	if p.Info.Org == "" {
		p.Info.Org = util.KeyUnknown
	}
	p.Info.Org = promptString("Enter the github organization that owns this project", p.Info.Org)

	if p.Package == "" || p.Package == "github.com//" {
		p.Package = "github.com/" + p.Info.Org + "/" + p.Key
	}
	p.Package = promptString("Enter your project's package", p.Package)

	ph := "https://" + p.Package
	if p.Info.Homepage == "" {
		p.Info.Homepage = ph
	}
	p.Info.Homepage = promptString("Enter this project's home page", p.Info.Homepage)

	if p.Info.Sourcecode == "" {
		p.Info.Sourcecode = ph
	}
	p.Info.Sourcecode = promptString("Enter this project's source repository", p.Info.Sourcecode)
}

var promptTotal = 0

func promptString(query string, curr string) string {
	promptTotal++
	clilog(fmt.Sprint(promptTotal) + ": ")
	clilog(query)
	if curr == "" {
		clilog(" (optional)")
	} else {
		clilog(" (default: " + curr + ")")
	}
	clilog(util.StringDefaultLinebreak)
	clilog(" > ")
	text, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		clilog("error: " + err.Error() + util.StringDefaultLinebreak)
	}
	return util.OrDefault(strings.TrimSuffix(strings.TrimSuffix(text, "\n"), "\r"), curr)
}
