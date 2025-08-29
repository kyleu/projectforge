package project

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

const DefaultGoVersion = "1.25.0"

type Doc struct {
	Name     string `json:"name"`
	Provider string `json:"provider,omitempty"`
	URL      string `json:"url"`
}

type Docs []*Doc

type ModuleDef struct {
	Key  string `json:"key"`
	Path string `json:"path"`
	URL  string `json:"url"`
}

type ModuleDefs []*ModuleDef

type Info struct {
	Org             string            `json:"org,omitempty"`
	AuthorID        string            `json:"authorID,omitempty"`
	AuthorName      string            `json:"authorName,omitempty"`
	AuthorEmail     string            `json:"authorEmail,omitempty"`
	Team            string            `json:"team,omitempty"`
	License         string            `json:"license,omitempty"`
	Homepage        string            `json:"homepage,omitempty"`
	Sourcecode      string            `json:"sourcecode,omitempty"`
	Summary         string            `json:"summary,omitempty"`
	Description     string            `json:"description,omitempty"`
	CI              string            `json:"ci,omitempty"`
	Homebrew        string            `json:"homebrew,omitempty"`
	Bundle          string            `json:"bundle,omitempty"`
	SigningIdentity string            `json:"signingIdentity,omitempty"`
	NotarizeEmail   string            `json:"notarizeEmail,omitempty"`
	DatabaseEngine  string            `json:"databaseEngine,omitempty"`
	Slack           string            `json:"slack,omitempty"`
	Channels        []string          `json:"channels,omitempty"`
	JavaPackage     string            `json:"javaPackage,omitempty"`
	GoVersion       string            `json:"goVersion,omitempty"`
	GoBinary        string            `json:"goBinary,omitempty"`
	ConfigVars      util.KeyTypeDescs `json:"configVars,omitempty"`
	ExtraFiles      []string          `json:"extraFiles,omitempty"`
	IgnoredFiles    []string          `json:"ignoredFiles,omitempty"`
	Deployments     []string          `json:"deployments,omitempty"`
	EnvVars         []string          `json:"envVars,omitempty"`
	DockerPackages  []string          `json:"dockerPackages,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitempty"`
	Docs            Docs              `json:"docs,omitempty"`
	Acronyms        []string          `json:"acronyms,omitempty"`
	ModuleDefs      ModuleDefs        `json:"moduleDefs,omitempty"`
}

func (i *Info) SigningIdentityTrimmed() string {
	if strings.Contains(i.SigningIdentity, "(") && strings.Contains(i.SigningIdentity, ")") {
		return i.SigningIdentity[strings.LastIndex(i.SigningIdentity, "(")+1 : strings.LastIndex(i.SigningIdentity, ")")]
	}
	return i.SigningIdentity
}

func (i *Info) AuthorIDSafe() string {
	if i.AuthorID == "" {
		if !strings.Contains(i.AuthorEmail, "@") {
			return "no_owner"
		}
		return i.AuthorEmail
	}
	spl := util.StringSplitAndTrim(i.AuthorID, " ")
	ret := util.NewStringSlice(make([]string, 0, len(spl)))
	lo.ForEach(spl, func(x string, _ int) {
		x = strings.ReplaceAll(x, ",", "")
		if !strings.HasPrefix(x, "@") {
			x = "@" + x
		}
		ret.Push(x)
	})
	return ret.Join(" ")
}

func (i *Info) NotarizationEmail() string {
	if i.NotarizeEmail != "" {
		return i.NotarizeEmail
	}
	return i.AuthorEmail
}
