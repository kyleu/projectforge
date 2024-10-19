package gql

import (
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

type GQLInfo struct {
	Org             string
	AuthorID        string
	AuthorName      string
	AuthorEmail     string
	Team            string
	License         string
	Homepage        string
	Sourcecode      string
	Summary         string
	Description     string
	CI              string
	Homebrew        string
	Bundle          string
	SigningIdentity string
	NotarizeEmail   string
	DatabaseEngine  string
	Slack           string
	Channels        []string
	JavaPackage     string
	GoVersion       string
	GoBinary        string
	ConfigVars      []string
	ExtraFiles      []string
	IgnoredFiles    []string
	Deployments     []string
	EnvVars         []string
	Docs            project.Docs
	Acronyms        []string
	ModuleDefs      project.ModuleDefs
}

func FromInfo(i *project.Info, logger util.Logger) *GQLInfo {
	return &GQLInfo{
		Org:             i.Org,
		AuthorID:        i.AuthorID,
		AuthorName:      i.AuthorName,
		AuthorEmail:     i.AuthorEmail,
		Team:            i.Team,
		License:         i.License,
		Homepage:        i.Homepage,
		Sourcecode:      i.Sourcecode,
		Summary:         i.Summary,
		Description:     i.Description,
		CI:              i.CI,
		Homebrew:        i.Homebrew,
		Bundle:          i.Bundle,
		SigningIdentity: i.SigningIdentity,
		NotarizeEmail:   i.NotarizeEmail,
		DatabaseEngine:  i.DatabaseEngine,
		Slack:           i.Slack,
		Channels:        i.Channels,
		JavaPackage:     i.JavaPackage,
		GoVersion:       i.GoVersion,
		GoBinary:        i.GoBinary,
		ConfigVars:      i.ConfigVars.Strings(),
		ExtraFiles:      i.ExtraFiles,
		IgnoredFiles:    i.IgnoredFiles,
		Deployments:     i.Deployments,
		EnvVars:         i.EnvVars,
		Docs:            i.Docs,
		Acronyms:        i.Acronyms,
		ModuleDefs:      i.ModuleDefs,
	}
}

type GQLProject struct {
	Key     string
	Name    string
	Icon    string
	Exec    string
	Version string
	Package string
	Args    string
	Port    int32
	Modules []string
	Ignore  []string
	Tags    []string

	Info  *GQLInfo
	Theme *theme.Theme
	Build *project.Build
	Files []string

	ExportArgs *model.Args
	Config     util.ValueMap
	Path       string
	Parent     string
	Error      string
}

func FromProject(p *project.Project, logger util.Logger) *GQLProject {
	return &GQLProject{
		Key:        p.Key,
		Name:       p.Name,
		Icon:       p.Icon,
		Exec:       p.Exec,
		Version:    p.Version,
		Package:    p.Package,
		Args:       p.Args,
		Port:       int32(p.Port),
		Modules:    p.Modules,
		Ignore:     p.Ignore,
		Tags:       p.Tags,
		Info:       FromInfo(p.Info, logger),
		Theme:      p.Theme,
		Build:      p.Build,
		Files:      p.Files,
		ExportArgs: p.ExportArgs,
		Config:     p.Config,
		Path:       p.Path,
		Parent:     p.Parent,
		Error:      p.Error,
	}
}
