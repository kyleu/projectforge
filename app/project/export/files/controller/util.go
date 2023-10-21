package controller

import (
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

const (
	incDel   = "cutil.QueryStringBool(rc, \"includeDeleted\")"
	incFalse = ", false"
)

func checkGrp(ret *golang.Block, grp *model.Column, override ...string) {
	if grp == nil {
		return
	}
	x := "ret"
	if len(override) > 0 {
		x = override[0]
	}
	ret.W("\t\tif %s.%s != %sArg {", x, grp.Proper(), grp.Camel())
	ret.W("\t\t\treturn \"\", errors.New(\"unauthorized: incorrect [%s]\")", grp.Camel())
	ret.W("\t\t}")
}

func controllerModelFromPath(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"FromPath", "func")
	ret.W("func %sFromPath(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (*%s, error) {", m.Package, m.ClassRef())
	pks := m.PKs()
	lo.ForEach(pks, func(col *model.Column, _ int) {
		controllerArgFor(col, ret, "nil", 1)
	})
	args := lo.Map(pks, func(x *model.Column, _ int) string {
		return x.Camel() + "Arg"
	})
	suffix := ""
	if m.IsSoftDelete() {
		suffix = ", includeDeleted"
		ret.W("\tincludeDeleted := rc.UserValue(\"includeDeleted\") != nil || " + incDel)
	}
	ret.W("\treturn as.Services.%s.Get(ps.Context, nil, %s%s, ps.Logger)", m.Proper(), strings.Join(args, ", "), suffix)
	ret.W("}")

	return ret
}
