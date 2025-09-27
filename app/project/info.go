package project

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

const DefaultGoVersion = "1.25.0"

type Doc struct {
	Name     string `json:"name"`
	Provider string `json:"provider,omitzero"`
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
	Org             string            `json:"org,omitzero"`
	AuthorID        string            `json:"authorID,omitzero"`
	AuthorName      string            `json:"authorName,omitzero"`
	AuthorEmail     string            `json:"authorEmail,omitzero"`
	Team            string            `json:"team,omitzero"`
	License         string            `json:"license,omitzero"`
	Homepage        string            `json:"homepage,omitzero"`
	Sourcecode      string            `json:"sourcecode,omitzero"`
	Summary         string            `json:"summary,omitzero"`
	Description     string            `json:"description,omitzero"`
	CI              string            `json:"ci,omitzero"`
	Homebrew        string            `json:"homebrew,omitzero"`
	Bundle          string            `json:"bundle,omitzero"`
	SigningIdentity string            `json:"signingIdentity,omitzero"`
	NotarizeEmail   string            `json:"notarizeEmail,omitzero"`
	DatabaseEngine  string            `json:"databaseEngine,omitzero"`
	Slack           string            `json:"slack,omitzero"`
	Channels        []string          `json:"channels,omitempty"`
	JavaPackage     string            `json:"javaPackage,omitzero"`
	GoVersion       string            `json:"goVersion,omitzero"`
	GoBinary        string            `json:"goBinary,omitzero"`
	ConfigVars      util.KeyTypeDescs `json:"configVars,omitzero"`
	AdditionalPorts map[string]int    `json:"additionalPorts,omitzero"`
	ExtraFiles      []string          `json:"extraFiles,omitempty"`
	IgnoredFiles    []string          `json:"ignoredFiles,omitempty"`
	Deployments     []string          `json:"deployments,omitempty"`
	EnvVars         []string          `json:"envVars,omitempty"`
	DockerPackages  []string          `json:"dockerPackages,omitempty"`
	Dependencies    map[string]string `json:"dependencies,omitzero"`
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
