package controller

import (
	"strings"

	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
)

func controllerRevision(m *model.Model) *golang.Block {
	hc := m.HistoryColumns(true)
	ret := golang.NewBlock(m.Camel()+"Detail", "func")
	ret.W("func %s%s(rc *fasthttp.RequestCtx) {", m.Proper(), hc.Col.Proper())
	ret.W("\tAct(\"%s.%s\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", m.Package, hc.Col.Name)
	ret.W("\t\tlatest, err := %sFromPath(rc, as, ps)", m.Package)
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\t%s, err := cutil.RCRequiredInt(rc, %q)", hc.Col.Camel(), hc.Col.Camel())
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	pkRefs := make([]string, 0, len(m.PKs()))
	for _, pk := range m.PKs() {
		pkRefs = append(pkRefs, "latest."+pk.Proper())
	}
	pkStuff := strings.Join(pkRefs, ", ")
	ret.W("\t\tret, err := as.Services.%s.Get%s(ps.Context, nil, %s, %s, ps.Logger)", m.Proper(), hc.Col.Proper(), pkStuff, hc.Col.Camel())
	ret.W("\t\tif err != nil {")
	ret.W("\t\t\treturn \"\", err")
	ret.W("\t\t}")
	ret.W("\t\tps.Title = ret.String()")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn Render(rc, as, &v%s.Detail{Model: ret}, ps, %q, ret.String())", m.Package, m.Package)
	ret.W("\t})")
	ret.W("}")
	return ret
}
