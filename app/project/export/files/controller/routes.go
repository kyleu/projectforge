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
	g.AddBlocks(routes(g, args))
	lo.ForEach(args.Models.SortedDisplay(), func(m *model.Model, _ int) {
		if len(m.Group) == 0 {
			g.AddImport(helper.ImpAppController)
		} else {
			g.AddImport(helper.AppImport("controller/" + m.GroupString("c", "")))
		}
	})
	return g.Render(addHeader, linebreak)
}

func routes(g *golang.File, args *model.Args) *golang.Block {
	ret := golang.NewBlock("routes", "func")
	ret.W("func generatedRoutes(r *mux.Router) {")
	rct := routeContent(args)
	if len(rct) > 0 {
		g.AddImport(helper.ImpHTTP)
	}
	lo.ForEach(rct, func(x string, _ int) {
		ret.W(x)
	})
	ret.W("}")
	return ret
}

func routeContent(args *model.Args) []string {
	out := make([]string, 0, 100)
	lo.ForEach(args.Models.SortedRoutes(), func(m *model.Model, _ int) {
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

	l := fmt.Sprintf("\tmakeRoute(r, http.MethodGet, \"/%s\", %s.%sList)", m.Route(), pkg, m.Proper())
	nf := fmt.Sprintf("\tmakeRoute(r, http.MethodGet, \"/%s/_new\", %s.%sCreateForm)", m.Route(), pkg, m.Proper())
	ns := fmt.Sprintf("\tmakeRoute(r, http.MethodPost, \"/%s/_new\", %s.%sCreate)", m.Route(), pkg, m.Proper())
	nr := fmt.Sprintf("\tmakeRoute(r, http.MethodGet, \"/%s/_random\", %s.%sRandom)", m.Route(), pkg, m.Proper())
	d := fmt.Sprintf("\tmakeRoute(r, http.MethodGet, \"/%s/%s\", %s.%sDetail)", m.Route(), pkn, pkg, m.Proper())
	ef := fmt.Sprintf("\tmakeRoute(r, http.MethodGet, \"/%s/%s/edit\", %s.%sEditForm)", m.Route(), pkn, pkg, m.Proper())
	es := fmt.Sprintf("\tmakeRoute(r, http.MethodPost, \"/%s/%s/edit\", %s.%sEdit)", m.Route(), pkn, pkg, m.Proper())
	dl := fmt.Sprintf("\tmakeRoute(r, http.MethodGet, \"/%s/%s/delete\", %s.%sDelete)", m.Route(), pkn, pkg, m.Proper())
	out = append(out, l, nf, ns, nr, d, ef, es, dl)
	return out
}
