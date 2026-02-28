package action

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"projectforge.dev/projectforge/app/module"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

func cliGather(_ context.Context, p *project.Project, mods module.Modules, _ util.Logger) error {
	if err := cliGatherBasics(p); err != nil {
		return err
	}

	if err := cliGatherInfo(p); err != nil {
		return err
	}

	if p.Port == 0 {
		p.Port = 20000
	}
	portStr, err := promptString("Port", "Enter the default port your http server will run on", fmt.Sprint(p.Port))
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
	p.Info.License, err = promptString("License", "Enter the license used by this project", p.Info.License)
	if err != nil {
		return err
	}

	return nil
}

func cliGatherBasics(p *project.Project) error {
	if p.Key == "" {
		path, _ := os.Getwd()
		_, path = util.StringSplitPath(path)
		p.Key = path
	}
	var err error
	p.Key, err = promptString("Key", "Enter a project key; must only contain alphanumerics", strings.ToLower(p.Key))
	if err != nil {
		return err
	}

	if p.Name == "" {
		p.Name = p.Key
	}
	p.Name, err = promptString("Name", "Enter a project name; use title case and spaces if needed", p.Name)
	if err != nil {
		return err
	}

	if p.Icon == "" {
		p.Icon = "star"
	}

	if p.Exec == "" {
		p.Exec = p.Key
	}

	p.Version, err = promptString("Version", "Enter a version, such as 0.0.0", p.Version)
	if err != nil {
		return err
	}

	return nil
}

func cliGatherInfo(p *project.Project) error {
	if p.Info.Org == "" {
		p.Info.Org = util.KeyUnknown
	}
	org, err := promptString("GitHub Org", "Enter the github organization that owns this project", p.Info.Org)
	if err != nil {
		return err
	}
	p.Info.Org = org

	if p.Package == "" || p.Package == "github.com//" {
		p.Package = "github.com/" + p.Info.Org + "/" + p.Key
	}
	pkg, err := promptString("Package", "Enter your project's package", p.Package)
	if err != nil {
		return err
	}
	p.Package = pkg

	ph := "https://" + p.Package
	if p.Info.Homepage == "" {
		p.Info.Homepage = ph
	}
	homepage, err := promptString("Homepage", "Enter this project's home page", p.Info.Homepage)
	if err != nil {
		return err
	}
	p.Info.Homepage = homepage

	if p.Info.Sourcecode == "" {
		p.Info.Sourcecode = ph
	}
	source, err := promptString("Source", "Enter this project's source repository", p.Info.Sourcecode)
	if err != nil {
		return err
	}
	p.Info.Sourcecode = source

	if p.Info.Summary == "" {
		p.Info.Summary = "A simple project"
	}
	p.Info.Summary, err = promptString("Summary", "Enter a one-line description of this project", p.Info.Summary)
	if err != nil {
		return err
	}

	authName := util.Choose(p.Info.AuthorName == "", p.Info.Org, p.Info.AuthorName)
	p.Info.AuthorName, err = promptString("Author Name", "Enter the name of this project's owner", authName)
	if err != nil {
		return err
	}

	if p.Info.AuthorEmail == "" {
		p.Info.AuthorEmail = fmt.Sprintf("dev@%s.com", p.Key)
	}
	p.Info.AuthorEmail, err = promptString("Author Email", "Enter the email address of this project's owner", p.Info.AuthorEmail)
	if err != nil {
		return err
	}

	const ghText = "Enter the GitHub username(s) of this project's owner"
	p.Info.AuthorID, err = promptString("Author GitHub", ghText, util.Choose(p.Info.AuthorID == "", p.Info.Org, p.Info.AuthorID))
	if err != nil {
		return err
	}
	p.Info.Team, err = promptString("Team", "Enter the team that owns this project", p.Info.Team)
	if err != nil {
		return err
	}

	return nil
}
