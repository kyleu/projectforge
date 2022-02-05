package controller

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

func controllerHistory(m *model.Model) *golang.Block {
	ret := golang.NewBlock(m.Camel()+"History", "func")
	ret.W("func %sHistory(rc *fasthttp.RequestCtx) {", m.Proper())
	ret.W("\tact(\"%s.history\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package)
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Camel())
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\thistID, err := RCRequiredUUID(rc, \"historyID\")")
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", errors.Wrap(err, \"must provide [historyID] as an argument\")")
	ret.W("\t\t}")
	ret.W("\t\thist, err := as.Services.%s.GetHistory(ps.Context, nil, *histID)", m.Proper())
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Title = hist.ID.String()")
	ret.W("\t\tps.Data = hist")
	ret.W("\t\treturn render(rc, as, &v%s.History{Model: ret, History: hist}, ps, %q, ret.String(), hist.ID.String())", m.Package, m.Package)
	ret.W("\t})")
	ret.W("}")
	return ret
}