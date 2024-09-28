package controller

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/golang"
)

const msgEqSPrint = "\t\tmsg := fmt.Sprintf(\""

func controllerCreateForm(m *model.Model, grp *model.Column, models model.Models, prefix string) *golang.Block {
	ret := blockFor(m, prefix, grp, "create", "form")
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
	}
	var decls []string
	if grp != nil {
		decls = append(decls, fmt.Sprintf("%s: %sArg", grp.Proper(), grp.Camel()))
	}
	ret.WF("\t\tret := &%s{%s}", m.ClassRef(), strings.Join(decls, ", "))
	ret.W("\t\tif r.URL.Query().Get(\"prototype\") == util.KeyRandom {")
	ret.WF("\t\t\tret = %s.Random()", m.Package)
	var encounteredRelTables []string
	for _, rel := range m.Relations {
		relModel := models.Get(rel.Table)
		if !relModel.CanTraverseRelation() {
			continue
		}
		srcCol := rel.SrcColumns(m)[0]
		tgtCol := rel.TgtColumns(relModel)[0]
		if !slices.Contains(encounteredRelTables, rel.Table) {
			ret.WF("\t\t\trandom%s, err := as.Services.%s.Random(ps.Context, nil, ps.Logger)", relModel.Proper(), relModel.Proper())
			encounteredRelTables = append(encounteredRelTables, rel.Table)
		}
		ret.WF("\t\t\tif err == nil && random%s != nil {", relModel.Proper())
		var ref string
		if srcCol.Nullable && !srcCol.Type.Scalar() && !tgtCol.Nullable {
			ref = "&"
		}
		ret.WF("\t\t\t\tret.%s = %srandom%s.%s", srcCol.Proper(), ref, relModel.Proper(), tgtCol.Proper())
		ret.W("\t\t\t}")
	}

	ret.W("\t\t}")

	if grp == nil {
		ret.WF("\t\tps.SetTitleAndData(\"Create [%s]\", ret)", m.Proper())
	} else {
		ret.WF("\t\tps.SetTitleAndData(fmt.Sprintf(\"Create ["+m.Proper()+"] for %s [%%%%s]\", %sArg), ret)", grp.TitleLower(), grp.Camel())
	}
	ret.W("\t\tps.Data = ret")
	ret.WF("\t\treturn %sRender(r, as, &v%s.Edit{Model: ret, IsNew: true}, ps, %s%s, \"Create\")", prefix, m.Package, m.Breadcrumbs(), grp.BC())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerRandom(m *model.Model, prefix string) *golang.Block {
	ret := blockFor(m, prefix, nil, "random")
	ret.WF("\t\tret, err := as.Services.%s.Random(ps.Context, nil, ps.Logger)", m.Proper())
	ret.W("\t\tif err != nil {")
	ret.WF("\t\t\treturn \"\", errors.Wrap(err, \"unable to find random %s\")", m.Proper())
	ret.W("\t\t}")
	ret.W("\t\treturn ret.WebPath(), nil")
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerCreate(m *model.Model, grp *model.Column, prefix string) *golang.Block {
	ret := blockFor(m, prefix, grp, "create")
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
	}
	ret.WF("\t\tret, err := %sFromForm(r, ps.RequestBody, true)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.WF("\t\t\treturn \"\", errors.Wrap(err, \"unable to parse %s from form\")", m.Proper())
	ret.W("\t\t}")
	checkGrp(ret, grp)
	ret.WF("\t\terr = as.Services.%s.Create(ps.Context, nil, ps.Logger, ret)", m.Proper())
	ret.W("\t\tif err != nil {")
	ret.WF("\t\t\treturn \"\", errors.Wrap(err, \"unable to save newly-created %s\")", m.Proper())
	ret.W("\t\t}")
	ret.W(msgEqSPrint + m.Proper() + " [%%s] created\", ret.String())")
	ret.WF("\t\treturn %sFlashAndRedir(true, msg, ret.WebPath(), ps)", prefix)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerEditForm(m *model.Model, grp *model.Column, prefix string) *golang.Block {
	ret := blockFor(m, prefix, grp, "edit", "form")
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
	}
	ret.WF("\t\tret, err := %sFromPath(r, as, ps)", m.Package)
	ret.WE(2, `""`)
	checkGrp(ret, grp)
	ret.W("\t\tps.SetTitleAndData(\"Edit \"+ret.String(), ret)")
	ret.WF("\t\treturn %sRender(r, as, &v%s.Edit{Model: ret}, ps, %s%s, ret.String())", prefix, m.Package, m.Breadcrumbs(), grp.BC())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerEdit(m *model.Model, grp *model.Column, prefix string) *golang.Block {
	ret := blockFor(m, prefix, grp, "edit")
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
	}
	ret.WF("\t\tret, err := %sFromPath(r, as, ps)", m.Package)
	ret.WE(2, `""`)
	checkGrp(ret, grp)
	ret.WF("\t\tfrm, err := %sFromForm(r, ps.RequestBody, false)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.WF("\t\t\treturn \"\", errors.Wrap(err, \"unable to parse %s from form\")", m.Proper())
	ret.W("\t\t}")
	checkGrp(ret, grp, "frm")
	lo.ForEach(m.PKs(), func(pk *model.Column, _ int) {
		ret.WF("\t\tfrm.%s = ret.%s", pk.Proper(), pk.Proper())
	})
	ret.WF("\t\terr = as.Services.%s.Update(ps.Context, nil, frm, ps.Logger)", m.Proper())
	ret.W("\t\tif err != nil {")
	ret.WF("\t\t\treturn \"\", errors.Wrapf(err, \"unable to update %s [%%%%s]\", frm.String())", m.Proper())
	ret.W("\t\t}")
	ret.W(msgEqSPrint + m.Proper() + " [%%s] updated\", frm.String())")
	ret.WF("\t\treturn %sFlashAndRedir(true, msg, frm.WebPath(), ps)", prefix)
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerDelete(m *model.Model, grp *model.Column, prefix string) *golang.Block {
	ret := blockFor(m, prefix, grp, "delete")
	if grp != nil {
		controllerArgFor(grp, ret, `""`, 2)
	}
	ret.WF("\t\tret, err := %sFromPath(r, as, ps)", m.Package)
	ret.WE(2, `""`)
	checkGrp(ret, grp)
	pkCamels := lo.Map(m.PKs(), func(pk *model.Column, _ int) string {
		return "ret." + pk.Proper()
	})
	ret.WF("\t\terr = as.Services.%s.Delete(ps.Context, nil, %s, ps.Logger)", m.Proper(), strings.Join(pkCamels, ", "))
	ret.W("\t\tif err != nil {")
	ret.WF("\t\t\treturn \"\", errors.Wrapf(err, \"unable to delete %s [%%%%s]\", ret.String())", m.TitleLower())
	ret.W("\t\t}")
	ret.W(msgEqSPrint + m.Proper() + " [%%s] deleted\", ret.String())")
	ret.WF("\t\treturn %sFlashAndRedir(true, msg, \"/%s\", ps)", prefix, m.Route())
	ret.W("\t})")
	ret.W("}")
	return ret
}

func controllerModelFromForm(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Package+"FromForm", "func")
	ret.WF("func %sFromForm(r *http.Request, b []byte, setPK bool) (*%s, error) {", m.Package, m.ClassRef())
	ret.W("\tfrm, err := cutil.ParseForm(r, b)")
	ret.W("\tif err != nil {")
	ret.WF("\t\treturn nil, errors.Wrap(err, \"unable to parse form\")")
	ret.W("\t}")
	ret.WF("\tret, _, err := %s.FromMap(frm, setPK)", m.Package)
	ret.W("\treturn ret, err")
	ret.W("}")
	return ret
}
