package controller

import (
	"fmt"
	"strings"

	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func controllerDetail(models model.Models, m *model.Model, grp *model.Column, prefix string) *golang.Block {
	rrels := models.ReverseRelations(m.Name)
	ret := blockFor(m, prefix, grp, "detail")
	grpHistory := ""
	if grp != nil {
		controllerArgFor(grp, ret, "\"\"", 2)
		grpHistory = fmt.Sprintf(", %q", grp.Camel())
	}
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	checkGrp(ret, grp)
	checkRev(ret, m)
	if m.IsHistory() {
		ret.W("\t\thist, err := as.Services.%s.GetHistories(ps.Context, nil, %s, ps.Logger)", m.Proper(), m.PKs().ToRefs("ret."))
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn \"\", err")
		ret.W("\t\t}")
	}
	ret.W("\t\tps.Title = ret.TitleString() + \" (%s)\"", m.Title())
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
		ret.W("\t\tincDel := cutil.QueryStringBool(rc, \"includeDeleted\")")
	}

	var argKeys []string
	var argVals []string
	argAdd := func(k string, v string) {
		argKeys = append(argKeys, k)
		argVals = append(argVals, v)
	}
	argAdd("Model", "ret")

	if m.IsRevision() || m.IsHistory() || len(rrels) > 0 {
		argAdd("Params", "ps.Params")
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
		ret.W("\t\t%sPrms := ps.Params.Get(%q, nil, ps.Logger).Sanitize(%q)", rm.Camel(), rm.Package, rm.Package)
		const msg = "\t\t%sBy%s, err := as.Services.%s.GetBy%s(ps.Context, nil, %s, %sPrms%s, ps.Logger)"
		ret.W(msg, rm.CamelPlural(), rNames, rm.Proper(), rNames, lCols.ToRefs("ret.", rCols...), rm.Camel(), delSuffix)
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
		ret.W("\t\treturn %sRender(rc, as, &v%s.Detail{%s}, ps, %s%s, ret.String())", prefix, m.Package, argStr, m.Breadcrumbs(), grpHistory)
	} else {
		ret.W("\t\treturn %sRender(rc, as, &v%s.Detail{", prefix, m.Package)
		keyPad := util.StringArrayMaxLength(argKeys) + 1
		for idx, k := range argKeys {
			ret.W("\t\t\t%s %s,", util.StringPad(k+":", keyPad), argVals[idx])
		}
		ret.W("\t\t}, ps, %s%s, ret.String())", m.Breadcrumbs(), grpHistory)
	}
	ret.W("\t})")
	ret.W("}")
	return ret
}
