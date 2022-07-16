package grpc

import (
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

type FileArgs struct {
	Class string
	Pkg   string
	CPkg  string
	API   string
	Grp   *model.Column
}

func (a FileArgs) APISuffix() string {
	if a.API == "*" {
		return ""
	}
	return "By" + util.StringToTitle(a.API)
}

func (a FileArgs) KeySuffix() string {
	if a.API == "*" {
		return ""
	}
	return "." + a.API
}

func (a FileArgs) GrpSuffix() string {
	if a.Grp == nil {
		return ""
	}
	return "By" + a.Grp.Proper()
}

func (a FileArgs) AddStaticCheck(ref string, ret *golang.Block, m *model.Model, grp *model.Column, act string) {
	if grp == nil {
		return
	}
	ret.W("\tif %s.%s != %s {", ref, grp.Proper(), grp.Camel())
	const msg = "\t\treturn nil, errors.Errorf(\"unauthorized - user acting on behalf of [%%%%s] cannot %s a %s for [%%%%s]\", %s, %s.%s)"
	ret.W(msg, act, m.Camel(), grp.Camel(), ref, grp.Proper())
	ret.W("\t}")
}

func grpcAddSection(b *golang.Block, key string, indent int) {
	ind := util.StringRepeat("\t", indent)
	b.W(ind+"// $PF_SECTION_START(%s)$", key)
	b.W(ind+"// $PF_SECTION_END(%s)$", key)
}

func idClauseFor(m *model.Model) (string, string) {
	if m.IsSoftDelete() {
		return "\tincludeDeleted, _ := provider.GetBool(p.R, p.TX, \"includeDeleted\")", ", includeDeleted"
	}
	return "", ""
}

func grpcParamsFromRequest(m *model.Model, args string, g *golang.File) (*golang.Block, error) {
	ret := golang.NewBlock("grpcParamsFromRequest", "func")
	pks := m.PKs()
	ret.W("func %sParamsFromRequest(%s) (%s, error) {", m.Camel(), args, strings.Join(pks.GoTypeKeys(m.Package), ", "))
	zeroVals := strings.Join(pks.ZeroVals(), ", ")
	for _, col := range pks {
		err := grpcArgFor(col, ret, zeroVals, g)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
	}
	ret.W("\treturn %s, nil", strings.Join(pks.CamelNames(), ", "))
	ret.W("}")
	return ret, nil
}

func grpcArgFor(col *model.Column, b *golang.Block, zeroVals string, g *golang.File) error {
	switch col.Type.Key() {
	case types.KeyBool:
		b.W("\t%s, err := provider.GetRequestBool(p.R, %q)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn %s, err", zeroVals)
		b.W("\t}")
	case types.KeyInt:
		b.W("\t%s, err := provider.GetRequestInt(p.R, %q)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn %s, err", zeroVals)
		b.W("\t}")
	case types.KeyFloat:
		b.W("\t%s, err := provider.GetRequestFloat(p.R, %q)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn %s, err", zeroVals)
		b.W("\t}")
	case types.KeyString:
		b.W("\t%s, err := provider.GetRequestString(p.R, %q)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn %s, err", zeroVals)
		b.W("\t}")
	case types.KeyUUID:
		g.AddImport(helper.ImpUUID)
		b.W("\t%sString, err := provider.GetRequestString(p.R, %q)", col.Camel(), col.Camel())
		b.W("\tif err != nil {")
		b.W("\t\treturn %s, err", zeroVals)
		b.W("\t}")
		b.W("\t%sParsed := util.UUIDFromString(%sString)", col.Camel(), col.Camel())
		b.W("\tif %sParsed == nil {", col.Camel())
		b.W("\t\treturn %s, errors.New(\"field [%s] must be a uuid\")", zeroVals, col.Camel())
		b.W("\t}")
		b.W("\t%s := *%sParsed", col.Camel(), col.Camel())
	default:
		return errors.Errorf("unhandled gRPC arg type [%s]", col.Type.String())
	}
	return nil
}
