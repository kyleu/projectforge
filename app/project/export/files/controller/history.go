package controller

import (
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func controllerHistory(m *model.Model, prefix string) *golang.Block {
	ret := golang.NewBlock(m.Camel()+"History", "func")
	ret.W("func %sHistory(rc *fasthttp.RequestCtx) {", m.Proper())
	ret.W("\t%sAct(\"%s.history\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", prefix, m.Package)
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.WE(2, `""`)
	ret.W("\t\thistID, err := cutil.RCRequiredUUID(rc, \"historyID\")")
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"must provide [historyID] as an argument\")")
	ret.W("\t\t}")
	ret.W("\t\thistModel, err := as.Services.%s.GetHistory(ps.Context, nil, *histID, ps.Logger)", m.Proper())
	ret.WE(2, `""`)
	ret.W("\t\tps.Title = histModel.ID.String()")
	ret.W("\t\tps.Data = histModel")
	msg := "\t\treturn %sRender(rc, as, &v%s.History{Model: ret, History: histModel}, ps, %s, ret.String(), histModel.ID.String())"
	ret.W(msg, prefix, m.Package, m.Breadcrumbs())
	ret.W("\t})")
	ret.W("}")
	return ret
}
