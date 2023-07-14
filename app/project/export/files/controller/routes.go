package controller

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

func Routes(args *model.Args, addHeader bool, linebreak string) (*file.File, error) {
	g := golang.NewFile("routes", []string{"app", "controller", "routes"}, "generated")
	g.AddImport(helper.ImpRouter)
	g.AddBlocks(routes(args))
	lo.ForEach(args.Models, func(m *model.Model, _ int) {
		if len(m.Group) == 0 {
			g.AddImport(helper.ImpAppController)
		} else {
			g.AddImport(helper.AppImport("app/controller/" + m.GroupString("c", "")))
		}
	})
	return g.Render(addHeader, linebreak)
}

func routes(args *model.Args) *golang.Block {
	ret := golang.NewBlock("routes", "func")
	complexity := lo.Sum(lo.Map(args.Models, func(m *model.Model, _ int) int {
		return (len(m.GroupedColumns()) * 7) + 8
	}))
	if complexity > 80 {
		ret.W("//nolint:funlen")
	}

	ret.W("func generatedRoutes(r *router.Router) {")
	rct := routeContent(args)
	lo.ForEach(rct, func(x string, _ int) {
		ret.W(x)
	})
	ret.W("}")
	return ret
}

func routeContent(args *model.Args) []string {
	out := make([]string, 0, 100)
	lo.ForEach(args.Models, func(m *model.Model, _ int) {
		out = append(out, routeModelContent(m)...)
	})
	return out
}

func routeModelContent(m *model.Model) []string {
	out := make([]string, 0, 100)
	pkNames := make([]string, 0, len(m.PKs()))
	lo.ForEach(m.PKs(), func(col *model.Column, _ int) {
		pkNames = append(pkNames, fmt.Sprintf("{%s}", col.Camel()))
	})
	pkn := strings.Join(pkNames, "/")

	pkg := "controller"
	if len(m.Group) > 0 {
		pkg = m.LastGroup("c", "")
	}

	lo.ForEach(m.GroupedColumns(), func(grp *model.Column, _ int) {
		pathExtra := fmt.Sprintf("/%s/{%s}", grp.Camel(), grp.Camel())
		callSuffix := fmt.Sprintf("By%s", grp.Proper())

		g := fmt.Sprintf("\tr.GET(\"/%s/%s\", %s.%s%sList)", m.Route(), grp.Camel(), pkg, m.Proper(), grp.Proper())
		l := fmt.Sprintf("\tr.GET(\"/%s%s\", %s.%sList%s)", m.Route(), pathExtra, pkg, m.Proper(), callSuffix)
		nf := fmt.Sprintf("\tr.GET(\"/%s%s/new\", %s.%sCreateForm%s)", m.Route(), pathExtra, pkg, m.Proper(), callSuffix)
		ns := fmt.Sprintf("\tr.POST(\"/%s%s/new\", %s.%sCreate%s)", m.Route(), pathExtra, pkg, m.Proper(), callSuffix)
		d := fmt.Sprintf("\tr.GET(\"/%s%s/%s\", %s.%sDetail%s)", m.Route(), pathExtra, pkn, pkg, m.Proper(), callSuffix)
		ef := fmt.Sprintf("\tr.GET(\"/%s%s/%s/edit\", %s.%sEditForm%s)", m.Route(), pathExtra, pkn, pkg, m.Proper(), callSuffix)
		es := fmt.Sprintf("\tr.POST(\"/%s%s/%s/edit\", %s.%sEdit%s)", m.Route(), pathExtra, pkn, pkg, m.Proper(), callSuffix)
		out = append(out, g, l, nf, ns, d, ef, es)
	})

	l := fmt.Sprintf("\tr.GET(\"/%s\", %s.%sList)", m.Route(), pkg, m.Proper())
	nr := fmt.Sprintf("\tr.GET(\"/%s/random\", %s.%sCreateFormRandom)", m.Route(), pkg, m.Proper())
	nf := fmt.Sprintf("\tr.GET(\"/%s/new\", %s.%sCreateForm)", m.Route(), pkg, m.Proper())
	ns := fmt.Sprintf("\tr.POST(\"/%s/new\", %s.%sCreate)", m.Route(), pkg, m.Proper())
	d := fmt.Sprintf("\tr.GET(\"/%s/%s\", %s.%sDetail)", m.Route(), pkn, pkg, m.Proper())
	ef := fmt.Sprintf("\tr.GET(\"/%s/%s/edit\", %s.%sEditForm)", m.Route(), pkn, pkg, m.Proper())
	es := fmt.Sprintf("\tr.POST(\"/%s/%s/edit\", %s.%sEdit)", m.Route(), pkn, pkg, m.Proper())
	dl := fmt.Sprintf("\tr.GET(\"/%s/%s/delete\", %s.%sDelete)", m.Route(), pkn, pkg, m.Proper())
	out = append(out, l, nr, nf, ns, d, ef, es, dl)
	if m.IsHistory() {
		const msg = "\tr.GET(\"/%s/%s/history/{historyID}\", %s.%sHistory)"
		out = append(out, fmt.Sprintf(msg, m.Route(), strings.Join(pkNames, "/"), pkg, m.Proper()))
	}
	if m.IsRevision() {
		rc := m.HistoryColumn()
		const msg = "\tr.GET(\"/%s/%s/%s/{%s}\", %s.%s%s)"
		out = append(out, fmt.Sprintf(msg, m.Route(), strings.Join(pkNames, "/"), rc.Name, rc.Name, pkg, m.Proper(), rc.Proper()))
	}
	return out
}
