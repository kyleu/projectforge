package controller

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
)

func Routes(args *model.Args, linebreak string) (*file.File, error) {
	g := golang.NewFile("routes", []string{"app", "controller", "routes"}, "generated")
	g.AddImport(helper.ImpRouter)
	g.AddBlocks(routeConsts(), routes(g, args))
	lo.ForEach(args.Models.WithController().SortedRoutes(), func(m *model.Model, _ int) {
		if len(m.Group) == 0 {
			g.AddImport(helper.ImpAppController)
		} else {
			g.AddImport(helper.AppImport("controller/" + m.GroupString("c", "")))
		}
		g.AddBlocks(routeModelContent(m))
	})
	return g.Render(linebreak)
}

func routeConsts() *golang.Block {
	ret := golang.NewBlock("routeConsts", "const")
	ret.W(`const routeNew, routeRandom, routeEdit, routeDelete = "/_new", "/_random", "/edit", "/delete"`)
	return ret
}

func routes(g *golang.File, args *model.Args) *golang.Block {
	ret := golang.NewBlock("generatedRoutes", "func")
	ret.W("func generatedRoutes(r *mux.Router) {")
	rct := routeContent(args)
	if len(rct) > 0 {
		g.AddImport(helper.ImpHTTP)
	}
	ret.WA(rct...)
	ret.W("}")
	return ret
}

func routeContent(args *model.Args) []string {
	out := make([]string, 0, 100)
	lo.ForEach(args.Models.WithRoutes().SortedRoutes(), func(m *model.Model, _ int) {
		out = append(out, routeModelLink(m)...)
	})
	return out
}

func routeModelLink(m *model.Model) []string {
	ret := fmt.Sprintf("\tgeneratedRoutes%s(r, %q)", m.Proper(), "/"+m.Route())
	return []string{ret}
}

func routeModelContent(m *model.Model) *golang.Block {
	ret := golang.NewBlock("generatedRoutes"+m.Proper(), "func")
	pkNames := make([]string, 0, len(m.PKs()))
	lo.ForEach(m.PKs(), func(col *model.Column, _ int) {
		pkNames = append(pkNames, fmt.Sprintf("{%s}", col.Camel()))
	})
	pkn := strings.Join(pkNames, "/")

	pkg := "controller"
	if len(m.Group) > 0 {
		pkg = m.LastGroup("c", "")
	}

	rt := func(method string, extra string, act string) {
		if extra == "" {
			ret.WF("\tmakeRoute(r, http.Method%s, prefix, %s.%s%s)", method, pkg, m.Proper(), act)
		} else {
			ret.WF("\tmakeRoute(r, http.Method%s, prefix+%s, %s.%s%s)", method, extra, pkg, m.Proper(), act)
		}
	}

	ret.WF("func generatedRoutes%s(r *mux.Router, prefix string) {", m.Proper())
	ret.WF("\tconst pkn = %q", "/"+pkn)
	rt("Get", "", "List")
	rt("Get", "routeNew", "CreateForm")
	rt("Post", "routeNew", "Create")
	rt("Get", "routeRandom", "Random")
	rt("Get", "pkn", "Detail")
	rt("Get", "pkn+routeEdit", "EditForm")
	rt("Post", "pkn+routeEdit", "Edit")
	rt("Get", "pkn+routeDelete", "Delete")
	ret.WF("}")

	return ret
}
