package controller

import (
	"strings"

	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
)

const incDel = "cutil.RequestCtxBool(rc, \"includeDeleted\")"

func checkRev(ret *golang.Block, m *model.Model) {
	if !m.IsRevision() {
		return
	}
	hc := m.HistoryColumn()

	prmsStr := m.PKs().ToRefs("ret.")
	const msg = "\t\t%s, err := as.Services.%s.GetAll%s(ps.Context, nil, %s, params.Get(%q, nil, ps.Logger).Sanitize(%q), false)"
	ret.W(msg, hc.CamelPlural(), m.Proper(), hc.ProperPlural(), prmsStr, m.Package, m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
}

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
	for _, col := range pks {
		controllerArgFor(col, ret, "nil", 1)
	}
	args := make([]string, 0, len(pks))
	for _, x := range pks {
		args = append(args, x.Camel()+"Arg")
	}
	suffix := ""
	if m.IsSoftDelete() {
		suffix = ", includeDeleted"
		ret.W("\tincludeDeleted := rc.UserValue(\"includeDeleted\") != nil || " + incDel)
	}
	ret.W("\treturn as.Services.%s.Get(ps.Context, nil, %s%s)", m.Proper(), strings.Join(args, ", "), suffix)
	ret.W("}")

	return ret
}
