package gql

import (
	"embed"
	"fmt"
	"path"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/lib/graphql"
	"projectforge.dev/projectforge/app/util"
)

//go:embed *
var FS embed.FS

type Schema struct {
	sch string
	svc *graphql.Service
	st  *app.State
	log util.Logger
}

func NewSchema(st *app.State, logger util.Logger) (*Schema, error) {
	sch, err := read("schema.graphql", 0)
	if err != nil {
		return nil, err
	}
	ret := &Schema{sch: sch, svc: st.GraphQL, st: st, log: logger}
	err = ret.svc.RegisterStringSchema(util.AppKey, util.AppName, ret.sch, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

const (
	includePrefix = "## include["
	includeSuffix = "]"
)

func read(fn string, depth int) (string, error) {
	b, err := FS.ReadFile(fn)
	if err != nil {
		return "", err
	}
	lines := util.StringSplitLines(string(b))
	converted := lo.Map(lines, func(l string, _ int) string {
		s, e := strings.Index(l, includePrefix), strings.Index(l, includeSuffix)
		if s == -1 {
			return l
		}
		tgt, _ := path.Split(fn)
		tgt = path.Join(tgt, l[s+len(includePrefix):e])
		child, err := read(tgt, depth+1)
		if err != nil {
			return fmt.Sprintf("ERROR: %+v", err)
		}

		return child
	})
	ret := strings.Join(converted, "\n")
	return ret, nil
}
