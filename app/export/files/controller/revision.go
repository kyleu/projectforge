package controller

import (
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

func controllerRevision(m *model.Model) *golang.Block {
	hc := m.HistoryColumns(true)
	ret := golang.NewBlock(m.Camel()+"Detail", "func")
	ret.W("func %s%s(rc *fasthttp.RequestCtx) {", m.PackageProper(), hc.Col.Proper())
	ret.W("\tact(\"%s.%s\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package, hc.Col.Name)
	ret.W("\t\tret, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Title = ret.String()")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn render(rc, as, &v%s.Detail{Model: ret}, ps, %q, ret.String())", m.Package, m.Package)
	ret.W("\t})")
	ret.W("}")
	return ret
}
