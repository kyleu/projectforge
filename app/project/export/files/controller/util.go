package controller

import (
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

const (
	incDel   = "cutil.QueryStringBool(ps.URI, \"includeDeleted\")"
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
	ret.WF("\t\tif %s.%s != %sArg {", x, grp.Proper(), grp.Camel())
	ret.WF("\t\t\treturn \"\", errors.New(\"unauthorized: incorrect [%s]\")", grp.Camel())
	ret.W("\t\t}")
}

func controllerModelFromPath(m *model.Model, enums enum.Enums) *golang.Block {
	ret := golang.NewBlock(m.Proper()+"FromPath", "func")
	ret.WF("func %sFromPath(r *http.Request, as *app.State, ps *cutil.PageState) (*%s, error) {", m.Package, m.ClassRef())
	pks := m.PKs()
	lo.ForEach(pks, func(col *model.Column, _ int) {
		controllerArgFor(col, ret, "nil", 1, enums)
	})
	args := lo.Map(pks, func(x *model.Column, _ int) string {
		if x.Nullable && !x.Type.Scalar() {
			return "&" + x.Camel() + helper.TextArg
		} else {
			return x.Camel() + helper.TextArg
		}
	})
	var suffix string
	if m.IsSoftDelete() {
		suffix = ", includeDeleted"
		ret.W("\tincludeDeleted := " + incDel)
	}
	ret.WF("\treturn as.Services.%s.Get(ps.Context, nil, %s%s, ps.Logger)", m.Proper(), util.StringJoin(args, ", "), suffix)
	ret.W("}")

	return ret
}
