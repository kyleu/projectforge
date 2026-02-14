package action

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/doctor"
	"projectforge.dev/projectforge/app/doctor/checks"
	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func clilog(s string) {
	print(s) //nolint:forbidigo
}

func cliProject(ctx context.Context, p *project.Project, mods module.Modules, logger util.Logger) error {
	clilog(util.AppName + "\nLet's create a new project!\n\n")
	var err error

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
	p.Key, err = promptString("Enter a project key; must only contain alphanumerics", strings.ToLower(p.Key))
	if err != nil {
		return err
	}

	if p.Name == "" {
		p.Name = p.Key
	}
	p.Name, err = promptString("Enter a project name; use title case and spaces if needed", p.Name)
	if err != nil {
		return err
	}

	if p.Icon == "" {
		p.Icon = "star"
	}

	if p.Exec == "" {
		p.Exec = p.Key
	}

	p.Version, err = promptString("Enter a version, such as 0.0.0", p.Version)
	if err != nil {
		return err
	}

	err = gatherProjectInfo(p)
	if err != nil {
		return err
	}

	if p.Info.Summary == "" {
		p.Info.Summary = "A simple project"
	}
	p.Info.Summary, err = promptString("Enter a one-line description of this project", p.Info.Summary)
	if err != nil {
		return err
	}

	authName := util.Choose(p.Info.AuthorName == "", p.Info.Org, p.Info.AuthorName)
	p.Info.AuthorName, err = promptString("Enter the name of this project's owner", authName)
	if err != nil {
		return err
	}

	if p.Info.AuthorEmail == "" {
		p.Info.AuthorEmail = fmt.Sprintf("dev@%s.com", p.Key)
	}
	p.Info.AuthorEmail, err = promptString("Enter the email address of this project's owner", p.Info.AuthorEmail)
	if err != nil {
		return err
	}

	p.Info.AuthorID, err = promptString("Enter the GitHub username(s) of this project's owner", util.Choose(p.Info.AuthorID == "", p.Info.Org, p.Info.AuthorID))
	if err != nil {
		return err
	}
	p.Info.Team, err = promptString("Enter the team that owns this project", p.Info.Team)
	if err != nil {
		return err
	}

	if p.Port == 0 {
		p.Port = 20000
	}
	portStr, err := promptString("Enter the default port your http server will run on", fmt.Sprint(p.Port))
	if err != nil {
		return err
	}
	prt, _ := strconv.ParseInt(portStr, 10, 32)
	p.Port = int(prt)

	p.Modules, err = promptModules(p.Modules, mods)
	if err != nil {
		return err
	}

	if p.Info.License == "" {
		p.Info.License = licenseProprietary
	}
	p.Info.License, err = promptString("Enter the license used by this project", p.Info.License)
	if err != nil {
		return err
	}

	return nil
}

func gatherProjectInfo(p *project.Project) error {
	if p.Info.Org == "" {
		p.Info.Org = util.KeyUnknown
	}
	org, err := promptString("Enter the github organization that owns this project", p.Info.Org)
	if err != nil {
		return err
	}
	p.Info.Org = org

	if p.Package == "" || p.Package == "github.com//" {
		p.Package = "github.com/" + p.Info.Org + "/" + p.Key
	}
	pkg, err := promptString("Enter your project's package", p.Package)
	if err != nil {
		return err
	}
	p.Package = pkg

	ph := "https://" + p.Package
	if p.Info.Homepage == "" {
		p.Info.Homepage = ph
	}
	homepage, err := promptString("Enter this project's home page", p.Info.Homepage)
	if err != nil {
		return err
	}
	p.Info.Homepage = homepage

	if p.Info.Sourcecode == "" {
		p.Info.Sourcecode = ph
	}
	source, err := promptString("Enter this project's source repository", p.Info.Sourcecode)
	if err != nil {
		return err
	}
	p.Info.Sourcecode = source
	return nil
}

var promptTotal = 0

func promptString(query string, curr string) (string, error) {
	promptTotal++
	text := curr
	in := huh.NewInput().Title(fmt.Sprintf("%d: %s", promptTotal, query)).Value(&text)
	if curr == "" {
		in = in.Description("Optional")
	} else {
		in = in.Description("Default: " + curr)
	}
	f := huh.NewForm(huh.NewGroup(in))
	err := f.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", errors.New("project creation canceled")
		}
		clilog("error: " + err.Error() + util.StringDefaultLinebreak)
		return "", err
	}
	return util.OrDefault(strings.TrimSpace(text), curr), nil
}
