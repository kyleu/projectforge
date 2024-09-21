package controller

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

func controllerDetail(g *golang.File, models model.Models, m *model.Model, grp *model.Column, audit bool, prefix string) *golang.Block {
	rrels := models.ReverseRelations(m.Name)
	ret := blockFor(m, prefix, grp, util.KeyDetail)
	grpHistory := ""
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
		grpHistory = fmt.Sprintf(", %q", grp.Camel())
	}
	ret.W("\t\tret, err := %sFromPath(r, as, ps)", m.Package)
	ret.WE(2, `""`)
	checkGrp(ret, grp)
	ret.W("\t\tps.SetTitleAndData(ret.TitleString()+\" (%s)\", ret)", m.Title())
	_, shouldIncDel := lo.Find(rrels, func(r *model.Relation) bool {
		return models.Get(r.Table).IsSoftDelete()
	})
	if shouldIncDel {
		ret.W("\t\tincDel := cutil.QueryStringBool(r, \"includeDeleted\")")
	}
	ret.WB()
	if audit {
		pks := m.PKs()
		if len(pks) > 1 {
			fields := lo.Map(pks, func(pk *model.Column, _ int) string {
				return fmt.Sprintf("%s: ret.%s", pk.Proper(), pk.Proper())
			})
			ret.W("\t\tpk := &%s.PK{%s}", m.Package, strings.Join(fields, ", "))
			msg := "\t\trelatedAuditRecords, err := as.Services.Audit.RecordsForModel(ps.Context, nil, %q, pk.String(), nil, ps.Logger)"
			ret.W(msg, m.Proper())
		} else {
			msg := "\t\trelatedAuditRecords, err := as.Services.Audit.RecordsForModel(ps.Context, nil, %q, %s, nil, ps.Logger)"
			ret.W(msg, m.Name, m.PKs().ToGoStrings("ret.", false, 160))
		}
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn \"\", errors.Wrapf(err, \"unable to retrieve related audit records\")")
		ret.W("\t\t}")
		ret.WB()
	}

	argKeys, argVals := getArgs(g, models, m, rrels, ret)
	revArgKeys, revArgVals := getReverseArgs(models, m, rrels, ret)
	if len(argKeys)+len(revArgKeys) <= 2 {
		args := lo.Map(argKeys, func(k string, idx int) string {
			return fmt.Sprintf("%s: %s", k, argVals[idx])
		})
		argStr := strings.Join(args, ", ")
		if audit {
			msg := "\t\treturn %sRender(r, as, &v%s.Detail{%s, AuditRecords: relatedAuditRecords}, ps, %s%s, %s)"
			ret.W(msg, prefix, m.Package, argStr, m.Breadcrumbs(), grpHistory, bcFor(m))
		} else {
			msg := "\t\treturn %sRender(r, as, &v%s.Detail{%s}, ps, %s%s, %s)"
			ret.W(msg, prefix, m.Package, argStr, m.Breadcrumbs(), grpHistory, bcFor(m))
		}
	} else {
		ret.W("\t\treturn %sRender(r, as, &v%s.Detail{", prefix, m.Package)
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
		ret.W("\t\t}, ps, %s%s, %s)", m.Breadcrumbs(), grpHistory, bcFor(m))
	}
	ret.W("\t})")
	ret.W("}")
	return ret
}

func bcFor(m *model.Model) any {
	icon := m.Icon
	var ret string
	if icons := m.Columns.WithFormat("icon"); len(icons) > 0 {
		ret = fmt.Sprintf("ret.TitleString()+\"**\"+ret.%s", icons[0].IconDerived())
	} else {
		ret = fmt.Sprintf("ret.TitleString()+\"**%s\"", icon)
	}
	if m.HasTag("menu-items") {
		if pks := m.PKs(); len(pks) == 1 {
			return "ret." + pks[0].Proper()
		}
	}
	return ret
}

func getArgs(g *golang.File, models model.Models, m *model.Model, rrels model.Relations, ret *golang.Block) ([]string, []string) {
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
			g.AddImport(helper.AppImport(rm.PackageWithGroup("")))
			ret.W("\t\tvar %sBy%s *%s.%s", rm.Camel(), lNames, rm.Package, rm.Proper())
			ret.W("\t\tif %s {", strings.Join(conditions, " && "))
			ret.W("\t\t\t%sBy%s, _ = as.Services.%s.Get(ps.Context, nil, %s%s, ps.Logger)", rm.Camel(), lNames, rm.Proper(), strings.Join(args, ", "), suffix)
			ret.W("\t\t}")
		}
	})
	if len(m.Relations) > 0 {
		ret.WB()
	}
	if len(rrels) > 0 {
		argAdd("Params", "ps.Params")
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
		ret.W("\t\trel%sBy%sPrms := ps.Params.Sanitized(%q, ps.Logger)", rm.ProperPlural(), rNames, rm.Package)
		const msg = "\t\trel%sBy%s, err := as.Services.%s.GetBy%s(ps.Context, nil, %s, rel%sBy%sPrms%s, ps.Logger)"
		ret.W(msg, rm.ProperPlural(), rNames, rm.Proper(), rNames, lCols.ToRefs("ret.", rCols...), rm.ProperPlural(), rNames, delSuffix)
		ret.W("\t\tif err != nil {")
		ret.W("\t\t\treturn \"\", errors.Wrap(err, \"unable to retrieve child %s\")", rm.TitlePluralLower())
		ret.W("\t\t}")
	})
	return argKeys, argVals
}
