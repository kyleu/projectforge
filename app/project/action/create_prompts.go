package action

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type CreatePromptKind string

const (
	CreatePromptString  = CreatePromptKind("string")
	CreatePromptModules = CreatePromptKind("modules")
)

const (
	createPromptKey         = "key"
	createPromptName        = "name"
	createPromptVersion     = "version"
	createPromptOrg         = "org"
	createPromptPackage     = "package"
	createPromptHomepage    = "homepage"
	createPromptSource      = "sourcecode"
	createPromptSummary     = "summary"
	createPromptAuthorName  = "authorName"
	createPromptAuthorEmail = "authorEmail"
	createPromptAuthorID    = "authorID"
	createPromptTeam        = "team"
	createPromptPort        = "port"
	createPromptModules     = "modules"
	createPromptLicense     = "license"
)

type CreatePrompt struct {
	Key   string
	Label string
	Query string
	Kind  CreatePromptKind
}

func CreatePrompts() []CreatePrompt {
	return []CreatePrompt{
		{Key: createPromptKey, Label: "Key", Query: "Enter a project key; must only contain alphanumerics", Kind: CreatePromptString},
		{Key: createPromptName, Label: "Name", Query: "Enter a project name; use title case and spaces if needed", Kind: CreatePromptString},
		{Key: createPromptVersion, Label: "Version", Query: "Enter a version, such as 0.0.0", Kind: CreatePromptString},
		{Key: createPromptOrg, Label: "GitHub Org", Query: "Enter the github organization that owns this project", Kind: CreatePromptString},
		{Key: createPromptPackage, Label: "Package", Query: "Enter your project's package", Kind: CreatePromptString},
		{Key: createPromptHomepage, Label: "Homepage", Query: "Enter this project's home page", Kind: CreatePromptString},
		{Key: createPromptSource, Label: "Source", Query: "Enter this project's source repository", Kind: CreatePromptString},
		{Key: createPromptSummary, Label: "Summary", Query: "Enter a one-line description of this project", Kind: CreatePromptString},
		{Key: createPromptAuthorName, Label: "Author Name", Query: "Enter the name of this project's owner", Kind: CreatePromptString},
		{Key: createPromptAuthorEmail, Label: "Author Email", Query: "Enter the email address of this project's owner", Kind: CreatePromptString},
		{Key: createPromptAuthorID, Label: "Author GitHub", Query: "Enter the GitHub username(s) of this project's owner", Kind: CreatePromptString},
		{Key: createPromptTeam, Label: "Team", Query: "Enter the team that owns this project", Kind: CreatePromptString},
		{Key: createPromptPort, Label: "Port", Query: "Enter the default port your http server will run on", Kind: CreatePromptString},
		{Key: createPromptModules, Label: "Modules", Query: "Select the modules your project will use", Kind: CreatePromptModules},
		{Key: createPromptLicense, Label: "License", Query: "Enter the license used by this project", Kind: CreatePromptString},
	}
}

func CreatePromptDefaultString(p *project.Project, promptKey string) string {
	if p == nil {
		return ""
	}
	p.Sanitize()
	switch promptKey {
	case createPromptKey:
		if p.Key == "" {
			path, _ := os.Getwd()
			_, path = util.StringSplitPath(path)
			p.Key = path
		}
		return strings.ToLower(p.Key)
	case createPromptName:
		if p.Name == "" {
			p.Name = p.Key
		}
		return p.Name
	case createPromptVersion:
		if p.Icon == "" {
			p.Icon = "star"
		}
		if p.Exec == "" {
			p.Exec = p.Key
		}
		return p.Version
	case createPromptOrg:
		if p.Info.Org == "" {
			p.Info.Org = util.KeyUnknown
		}
		return p.Info.Org
	case createPromptPackage:
		if p.Package == "" || p.Package == "github.com//" {
			p.Package = "github.com/" + p.Info.Org + "/" + p.Key
		}
		return p.Package
	case createPromptHomepage:
		ph := "https://" + p.Package
		if p.Info.Homepage == "" {
			p.Info.Homepage = ph
		}
		return p.Info.Homepage
	case createPromptSource:
		ph := "https://" + p.Package
		if p.Info.Sourcecode == "" {
			p.Info.Sourcecode = ph
		}
		return p.Info.Sourcecode
	case createPromptSummary:
		if p.Info.Summary == "" {
			p.Info.Summary = "A simple project"
		}
		return p.Info.Summary
	case createPromptAuthorName:
		return util.Choose(p.Info.AuthorName == "", p.Info.Org, p.Info.AuthorName)
	case createPromptAuthorEmail:
		if p.Info.AuthorEmail == "" {
			p.Info.AuthorEmail = fmt.Sprintf("dev@%s.com", p.Key)
		}
		return p.Info.AuthorEmail
	case createPromptAuthorID:
		return util.Choose(p.Info.AuthorID == "", p.Info.Org, p.Info.AuthorID)
	case createPromptTeam:
		return p.Info.Team
	case createPromptPort:
		if p.Port == 0 {
			p.Port = 20000
		}
		return fmt.Sprint(p.Port)
	case createPromptLicense:
		if p.Info.License == "" {
			p.Info.License = licenseProprietary
		}
		return p.Info.License
	default:
		return ""
	}
}

func CreatePromptDefaultModules(p *project.Project) []string {
	if p == nil {
		return nil
	}
	return append([]string{}, p.Modules...)
}

func ApplyCreatePromptString(p *project.Project, promptKey string, value string) error {
	if p == nil {
		return nil
	}
	p.Sanitize()
	switch promptKey {
	case createPromptKey:
		p.Key = value
	case createPromptName:
		p.Name = value
	case createPromptVersion:
		p.Version = value
	case createPromptOrg:
		p.Info.Org = value
	case createPromptPackage:
		p.Package = value
	case createPromptHomepage:
		p.Info.Homepage = value
	case createPromptSource:
		p.Info.Sourcecode = value
	case createPromptSummary:
		p.Info.Summary = value
	case createPromptAuthorName:
		p.Info.AuthorName = value
	case createPromptAuthorEmail:
		p.Info.AuthorEmail = value
	case createPromptAuthorID:
		p.Info.AuthorID = value
	case createPromptTeam:
		p.Info.Team = value
	case createPromptPort:
		prt, _ := strconv.ParseInt(value, 10, 32)
		p.Port = int(prt)
	case createPromptLicense:
		p.Info.License = value
	}
	return nil
}

func ApplyCreatePromptModules(p *project.Project, value []string) {
	if p == nil {
		return
	}
	p.Modules = util.ArraySorted(util.ArrayRemoveDuplicates(value))
}

func CreateConfigFromProject(p *project.Project) util.ValueMap {
	if p == nil {
		return util.ValueMap{}
	}
	p.Sanitize()
	ret := p.ToMap()
	ret["org"] = p.Info.Org
	ret["homepage"] = p.Info.Homepage
	ret["sourcecode"] = p.Info.Sourcecode
	ret["summary"] = p.Info.Summary
	ret["authorName"] = p.Info.AuthorName
	ret["authorEmail"] = p.Info.AuthorEmail
	ret["authorID"] = p.Info.AuthorID
	ret["team"] = p.Info.Team
	ret["license"] = p.Info.License
	ret["path"] = p.Path
	return ret
}

func CreateSummaryLines(p *project.Project) []string {
	if p == nil {
		return nil
	}
	p.Sanitize()
	return []string{
		fmt.Sprintf("Key: %s", p.Key),
		fmt.Sprintf("Name: %s", p.Name),
		fmt.Sprintf("Version: %s", p.Version),
		fmt.Sprintf("GitHub Org: %s", p.Info.Org),
		fmt.Sprintf("Package: %s", p.Package),
		fmt.Sprintf("Homepage: %s", p.Info.Homepage),
		fmt.Sprintf("Source: %s", p.Info.Sourcecode),
		fmt.Sprintf("Summary: %s", p.Info.Summary),
		fmt.Sprintf("Author Name: %s", p.Info.AuthorName),
		fmt.Sprintf("Author Email: %s", p.Info.AuthorEmail),
		fmt.Sprintf("Author GitHub: %s", p.Info.AuthorID),
		fmt.Sprintf("Team: %s", p.Info.Team),
		fmt.Sprintf("Port: %d", p.Port),
		fmt.Sprintf("Modules: %s", util.StringJoin(p.Modules, ", ")),
		fmt.Sprintf("License: %s", p.Info.License),
	}
}
