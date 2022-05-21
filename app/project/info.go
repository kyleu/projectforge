package project

import (
	"strings"

	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/util"
)

const defaultGoVersion = "1.18"

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
	License         string            `json:"license,omitempty"`
	Homepage        string            `json:"homepage,omitempty"`
	Sourcecode      string            `json:"sourcecode,omitempty"`
	Summary         string            `json:"summary,omitempty"`
	Description     string            `json:"description,omitempty"`
	CI              string            `json:"ci,omitempty"`
	Homebrew        string            `json:"homebrew,omitempty"`
	Bundle          string            `json:"bundle,omitempty"`
	SigningIdentity string            `json:"signingIdentity,omitempty"`
	Slack           string            `json:"slack,omitempty"`
	JavaPackage     string            `json:"javaPackage,omitempty"`
	GoVersion       string            `json:"goVersion,omitempty"`
	GoBinary        string            `json:"goBinary,omitempty"`
	ConfigVars      util.KeyTypeDescs `json:"configVars,omitempty"`
	ExtraFiles      []string          `json:"extraFiles,omitempty"`
	Deployments     []string          `json:"deployments,omitempty"`
	ModuleDefs      ModuleDefs        `json:"moduleDefs,omitempty"`
	ModuleArgs      util.ValueMap     `json:"moduleArgs,omitempty"`
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
	ret := make([]string, 0, len(spl))
	for _, x := range spl {
		x = strings.ReplaceAll(x, ",", "")
		if !strings.HasPrefix(x, "@") {
			x = "@" + x
		}
		ret = append(ret, x)
	}
	return strings.Join(ret, " ")
}

func (i *Info) ModuleArg(mod string) any {
	if i == nil || i.ModuleArgs == nil {
		return nil
	}
	return i.ModuleArgs[mod]
}

func (i *Info) ModuleArgExport() (*model.Args, error) {
	arg := i.ModuleArg("export")
	ret := &model.Args{}
	err := util.CycleJSON(arg, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
