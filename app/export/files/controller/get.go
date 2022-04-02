package controller

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/util"
)

const incDel = "cutil.RequestCtxBool(rc, \"includeDeleted\")"

func controllerList(m *model.Model, grp *model.Column) *golang.Block {
	ret := blockFor(m, grp, "list")
	meth := "List"
	grpArgs := ""
	if grp != nil {
		meth = "GetBy" + grp.Title()
		grpArgs = ", " + grp.Camel() + "Arg"
		controllerArgFor(grp, ret, "\"\"", 2)
	}

	ret.W("\t\tps.Title = %sDefaultTitle", m.Camel())
	ret.W("\t\tparams := cutil.ParamSetFromRequest(rc)")
	suffix := ""
	if m.IsSoftDelete() {
		suffix = ", " + incDel
	}
	ret.W("\t\tprms := params.Get(%q, nil, ps.Logger).Sanitize(%q)", m.Package, m.Package)
	const msg = "\t\tret, err := as.Services.%s.%s(ps.Context, nil%s, prms%s)"
	ret.W(msg, m.Proper(), meth, grpArgs, suffix)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn render(rc, as, &v%s.List{Models: ret, Params: params}, ps, %q%s)", m.Package, m.Package, grp.BC())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerDetail(models model.Models, m *model.Model, grp *model.Column) *golang.Block {
	rrels := models.ReverseRelations(m.Name)
	ret := blockFor(m, grp, "detail")
	grpHistory := ""
	if grp != nil {
		controllerArgFor(grp, ret, "\"\"", 2)
		grpHistory = fmt.Sprintf(", %q", grp.Camel())
	}
	if m.IsRevision() || m.IsHistory() || len(rrels) > 0 {
		ret.W("\t\tparams := cutil.ParamSetFromRequest(rc)")
	}
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	checkGrp(ret, grp)
	checkRev(ret, m)
	if m.IsHistory() {
		ret.W("\t\thist, err := as.Services.%s.GetHistories(ps.Context, nil, %s)", m.Proper(), m.PKs().ToRefs("ret."))
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn \"\", err")
		ret.W("\t\t}")
	}
	ret.W("\t\tps.Title = ret.String()")
	ret.W("\t\tps.Data = ret")

	var shouldIncDel bool
	for _, r := range rrels {
		rm := models.Get(r.Table)
		if rm.IsSoftDelete() {
			shouldIncDel = true
			break
		}
	}

	if shouldIncDel {
		ret.W("\t\tincDel := cutil.RequestCtxBool(rc, \"includeDeleted\")")
	}

	var argKeys []string
	var argVals []string
	argAdd := func(k string, v string) {
		argKeys = append(argKeys, k)
		argVals = append(argVals, v)
	}
	argAdd("Model", "ret")

	if m.IsRevision() || m.IsHistory() || len(rrels) > 0 {
		argAdd("Params", "params")
	}
	for _, rel := range rrels {
		rm := models.Get(rel.Table)
		delSuffix := ""
		if rm.IsSoftDelete() {
			delSuffix = ", incDel"
		}
		lCols := rel.SrcColumns(m)
		rCols := rel.TgtColumns(rm)
		rNames := strings.Join(rCols.ProperNames(), "")
		ret.W("\t\t%sPrms := params.Get(%q, nil, ps.Logger).Sanitize(%q)", rm.Camel(), rm.Camel(), rm.Camel())
		const msg = "\t\t%sBy%s, err := as.Services.%s.GetBy%s(ps.Context, nil, %s, %sPrms%s)"
		ret.W(msg, rm.CamelPlural(), rNames, rm.Proper(), rNames, lCols.ToRefs("ret."), rm.Camel(), delSuffix)
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to retrieve child %s\")", rm.TitlePluralLower())
		ret.W("\t\t}")
		argAdd(fmt.Sprintf("%sBy%s", rm.ProperPlural(), rNames), fmt.Sprintf("%sBy%s", rm.CamelPlural(), rNames))
	}
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		argAdd(revCol.ProperPlural(), revCol.CamelPlural())
	}
	if m.IsHistory() {
		argAdd("Histories", "hist")
	}
	if len(argKeys) <= 2 {
		args := make([]string, 0, len(argKeys))
		for idx, k := range argKeys {
			args = append(args, fmt.Sprintf("%s: %s", k, argVals[idx]))
		}
		argStr := strings.Join(args, ", ")
		ret.W("\t\treturn render(rc, as, &v%s.Detail{%s}, ps, %q%s, ret.String())", m.Package, argStr, m.Package, grpHistory)
	} else {
		ret.W("\t\treturn render(rc, as, &v%s.Detail{", m.Package)
		keyPad := util.StringArrayMaxLength(argKeys) + 1
		for idx, k := range argKeys {
			ret.W("\t\t\t%s %s,", util.StringPad(k+":", keyPad), argVals[idx])
		}
		ret.W("\t\t}, ps, %q%s, ret.String())", m.Package, grpHistory)
	}
	ret.W("\t})")
	ret.W("}")
	return ret
}

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
