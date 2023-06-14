package controller

import (
	"strings"

	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func controllerRevision(m *model.Model, prefix string) *golang.Block {
	hc := m.HistoryColumns(true)
	ret := golang.NewBlock(m.Camel()+"Detail", "func")
	ret.W("func %s%s(rc *fasthttp.RequestCtx) {", m.Proper(), hc.Col.Proper())
	ret.W("\t%sAct(\"%s.%s\", rc, func(as *app.State, ps *cutil.PageState) (string, error) {", prefix, m.Package, hc.Col.Name)
	ret.W("\t\tlatest, err := %sFromPath(rc, as, ps)", m.Package)
	ret.WE(2, `""`)
	ret.W("\t\t%s, err := cutil.RCRequiredInt(rc, %q)", hc.Col.Camel(), hc.Col.Camel())
	ret.WE(2, `""`)
	pkRefs := make([]string, 0, len(m.PKs()))
	for _, pk := range m.PKs() {
		pkRefs = append(pkRefs, "latest."+pk.Proper())
	}
	pkStuff := strings.Join(pkRefs, ", ")
	ret.W("\t\tret, err := as.Services.%s.Get%s(ps.Context, nil, %s, %s, ps.Logger)", m.Proper(), hc.Col.Proper(), pkStuff, hc.Col.Camel())
	ret.WE(2, `""`)
	ret.W("\t\tps.Title = ret.String()")
	ret.W("\t\tps.Data = ret")
	ret.W("\t\treturn %sRender(rc, as, &v%s.Detail{Model: ret}, ps, %s, ret.String())", prefix, m.Package, m.Breadcrumbs())
	ret.W("\t})")
	ret.W("}")
	return ret
}
