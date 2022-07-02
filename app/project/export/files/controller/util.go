package controller

import (
	"strings"

	"projectforge.dev/projectforge/app/project/export/golang"
	model2 "projectforge.dev/projectforge/app/project/export/model"
)

const (
	incDel   = "cutil.QueryStringBool(rc, \"includeDeleted\")"
	incFalse = ", false"
)

func checkRev(ret *golang.Block, m *model2.Model) {
	if !m.IsRevision() {
		return
	}
	hc := m.HistoryColumn()
	incDel := ""
	if m.IsSoftDelete() {
		incDel = incFalse
	}
	prmsStr := m.PKs().ToRefs("ret.")
	ret.W("\t\tprms := params.Get(%q, nil, ps.Logger).Sanitize(%q)", m.Package, m.Package)
	const msg = "\t\t%s, err := as.Services.%s.GetAll%s(ps.Context, nil, %s, prms%s, ps.Logger)"
	ret.W(msg, hc.CamelPlural(), m.Proper(), hc.ProperPlural(), prmsStr, incDel)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
}

func checkGrp(ret *golang.Block, grp *model2.Column, override ...string) {
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

func controllerModelFromPath(m *model2.Model) *golang.Block {
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
	ret.W("\treturn as.Services.%s.Get(ps.Context, nil, %s%s, ps.Logger)", m.Proper(), strings.Join(args, ", "), suffix)
	ret.W("}")

	return ret
}
