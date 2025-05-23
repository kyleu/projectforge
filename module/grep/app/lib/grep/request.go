package grep

import (
	"fmt"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Request struct {
	Path          string   `json:"path"`
	Query         string   `json:"query"`
	CaseSensitive bool     `json:"caseSensitive,omitempty"`
	Ignore        []string `json:"ignore,omitempty"`
}

func NewRequest(pth string, query string, caseSensitive bool, ignored ...string) *Request {
	return &Request{Path: pth, Query: query, CaseSensitive: caseSensitive, Ignore: ignored}
}

var defaultArgs = []string{"json", "max-columns=160", "max-columns-preview", "max-filesize=5M", "no-require-git", "sort=path", "stats", "with-filename"}

func (r *Request) ToCommand() string {
	if r.Path == "" {
		r.Path = "."
	}
	args := util.ArrayCopy(defaultArgs)
	if !r.CaseSensitive {
		args = append(args, "ignore-case")
	}
	for _, ign := range r.Ignore {
		args = append(args, "glob=!"+ign)
	}
	argString := util.StringJoin(lo.Map(util.ArraySorted(args), func(x string, _ int) string {
		return "--" + x
	}), " ")
	return fmt.Sprintf("rg %s %q, %q", argString, r.Query, r.Path)
}
