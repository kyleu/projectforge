// $PF_GENERATE_ONCE$
package gql

import (
	"embed"
	"fmt"
	"path"
	"strings"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/graphql"
	"{{{ .Package }}}/app/util"
)

//go:embed *
var FS embed.FS

type Schema struct {
	svc *graphql.Service
	sch string
}

func NewSchema(svc *graphql.Service) (*Schema, error) {
	sch, err := read("schema.graphql", 0)
	if err != nil {
		return nil, err
	}
	ret := &Schema{svc: svc, sch: sch}
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
