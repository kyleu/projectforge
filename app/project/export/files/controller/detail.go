package controller

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func controllerDetail(models model.Models, m *model.Model, grp *model.Column, audit bool, g *golang.File, prefix string) *golang.Block {
	rrels := models.ReverseRelations(m.Name)
	ret := blockFor(m, prefix, grp, (len(rrels)*6)+40, "detail")
	grpHistory := ""
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
		grpHistory = fmt.Sprintf(", %q", grp.Camel())
	}
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.WE(2, `""`)
	checkGrp(ret, grp)
	checkRev(ret, m)
	if m.IsHistory() {
		ret.W("\t\thist, err := as.Services.%s.GetHistories(ps.Context, nil, %s, ps.Logger)", m.Proper(), m.PKs().ToRefs("ret."))
		ret.WE(2, `""`)
	}
	ret.W("\t\tps.Title = ret.TitleString() + \" (%s)\"", m.Title())
	ret.W("\t\tps.Data = ret")

	_, shouldIncDel := lo.Find(rrels, func(r *model.Relation) bool {
		return models.Get(r.Table).IsSoftDelete()
	})
	if shouldIncDel {
		ret.W("\t\tincDel := cutil.QueryStringBool(rc, \"includeDeleted\")")
	}
	ret.WB()
	argKeys, argVals := getArgs(models, m, rrels, g, ret)
	revArgKeys, revArgVals := getReverseArgs(models, m, rrels, ret)
	if audit {
		ret.WB()
		msg := "\t\trelatedAuditRecords, err := as.Services.Audit.RecordsForModel(ps.Context, nil, %q, %s, nil, ps.Logger)"
		ret.W(msg, m.Name, m.PKs().ToGoStrings("ret."))
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn \"\", errors.Wrapf(err, \"unable to retrieve related audit records\")")
		ret.W("\t\t}")
		ret.WB()
	}
	if len(argKeys)+len(revArgKeys) <= 2 {
		args := lo.Map(argKeys, func(k string, idx int) string {
			return fmt.Sprintf("%s: %s", k, argVals[idx])
		})
		argStr := strings.Join(args, ", ")
		if audit {
			msg := "\t\treturn %sRender(rc, as, &v%s.Detail{%s, AuditRecords: relatedAuditRecords}, ps, %s%s, ret.String())"
			ret.W(msg, prefix, m.Package, argStr, m.Breadcrumbs(), grpHistory)
		} else {
			msg := "\t\treturn %sRender(rc, as, &v%s.Detail{%s}, ps, %s%s, ret.String())"
			ret.W(msg, prefix, m.Package, argStr, m.Breadcrumbs(), grpHistory)
		}
	} else {
		ret.W("\t\treturn %sRender(rc, as, &v%s.Detail{", prefix, m.Package)
		keyPad := util.StringArrayMaxLength(argKeys) + 1
		lo.ForEach(argKeys, func(k string, idx int) {
			ret.W("\t\t\t%s %s,", util.StringPad(k+":", keyPad), argVals[idx])
		})
		if len(revArgKeys) > 0 {
			revKeyPad := util.StringArrayMaxLength(revArgKeys) + 1
			ret.WB()
			lo.ForEach(revArgKeys, func(k string, idx int) {
				ret.W("\t\t\t%s %s,", util.StringPad(k+":", revKeyPad), revArgVals[idx])
			})
		}
		if audit {
			ret.WB()
			ret.W("\t\t\tAuditRecords: relatedAuditRecords,")
		}
		ret.W("\t\t}, ps, %s%s, ret.String())", m.Breadcrumbs(), grpHistory)
	}
	ret.W("\t})")
	ret.W("}")
	return ret
}

func getArgs(models model.Models, m *model.Model, rrels model.Relations, g *golang.File, ret *golang.Block) ([]string, []string) {
	var argKeys []string
	var argVals []string
	argAdd := func(k string, v string) {
		argKeys = append(argKeys, k)
		argVals = append(argVals, v)
	}
	argAdd("Model", "ret")
	lo.ForEach(m.Relations, func(rel *model.Relation, _ int) {
		rm := models.Get(rel.Table)
		lCols := rel.SrcColumns(m)

		lNames := strings.Join(lCols.ProperNames(), "")
		argAdd(fmt.Sprintf("%sBy%s", rm.Proper(), lNames), fmt.Sprintf("%sBy%s", rm.Camel(), lNames))

		var conditions []string
		var args []string
		lo.ForEach(lCols, func(col *model.Column, _ int) {
			if col.Nullable {
				conditions = append(conditions, fmt.Sprintf("ret.%s != nil", col.Proper()))
				args = append(args, fmt.Sprintf("*ret.%s", col.Proper()))
			} else {
				args = append(args, fmt.Sprintf("ret.%s", col.Proper()))
			}
		})
		suffix := rm.SoftDeleteSuffix()
		if len(conditions) == 0 {
			ret.W("\t\t%sBy%s, _ := as.Services.%s.Get(ps.Context, nil, %s%s, ps.Logger)", rm.Camel(), lNames, rm.Proper(), strings.Join(args, ", "), suffix)
		} else {
			g.AddImport(helper.AppImport("app/" + rm.PackageWithGroup("")))
			ret.W("\t\tvar %sBy%s *%s.%s", rm.Camel(), lNames, rm.Package, rm.Proper())
			ret.W("\t\tif %s {", strings.Join(conditions, " && "))
			ret.W("\t\t\t%sBy%s, _ = as.Services.%s.Get(ps.Context, nil, %s%s, ps.Logger)", rm.Camel(), lNames, rm.Proper(), strings.Join(args, ", "), suffix)
			ret.W("\t\t}")
		}
	})
	if len(m.Relations) > 0 {
		ret.WB()
	}
	if m.IsRevision() || m.IsHistory() || len(rrels) > 0 {
		argAdd("Params", "ps.Params")
	}
	if m.IsRevision() {
		revCol := m.HistoryColumn()
		argAdd(revCol.ProperPlural(), revCol.CamelPlural())
	}
	if m.IsHistory() {
		argAdd("Histories", "hist")
	}
	return argKeys, argVals
}

func getReverseArgs(models model.Models, m *model.Model, rrels model.Relations, ret *golang.Block) ([]string, []string) {
	argKeys := make([]string, 0, len(rrels))
	argVals := make([]string, 0, len(rrels))
	lo.ForEach(rrels, func(rrel *model.Relation, _ int) {
		rm := models.Get(rrel.Table)
		delSuffix := ""
		if rm.IsSoftDelete() {
			delSuffix = ", incDel"
		}
		lCols := rrel.SrcColumns(m)
		rCols := rrel.TgtColumns(rm)
		rNames := strings.Join(rCols.ProperNames(), "")
		argKeys = append(argKeys, fmt.Sprintf("Rel%sBy%s", rm.ProperPlural(), rNames))
		argVals = append(argVals, fmt.Sprintf("rel%sBy%s", rm.ProperPlural(), rNames))
		ret.W("\t\trel%sBy%sPrms := ps.Params.Get(%q, nil, ps.Logger).Sanitize(%q)", rm.ProperPlural(), rNames, rm.Package, rm.Package)
		const msg = "\t\trel%sBy%s, err := as.Services.%s.GetBy%s(ps.Context, nil, %s, rel%sBy%sPrms%s, ps.Logger)"
		ret.W(msg, rm.ProperPlural(), rNames, rm.Proper(), rNames, lCols.ToRefs("ret.", rCols...), rm.ProperPlural(), rNames, delSuffix)
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to retrieve child %s\")", rm.TitlePluralLower())
		ret.W("\t\t}")
	})
	return argKeys, argVals
}
