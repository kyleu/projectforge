package grpc

import (
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/export/files/helper"
	"projectforge.dev/projectforge/app/export/golang"
	"projectforge.dev/projectforge/app/export/model"
)

func grpcList(m *model.Model, grpcArgs string, grpcRet string, ga *FileArgs) *golang.Block {
	ret := golang.NewBlock("grpcList", "func")
	ret.W("func %sList%s(%s) %s {", m.Proper(), ga.APISuffix(), grpcArgs, grpcRet)
	idClause, suffix := idClauseFor(m)
	if idClause != "" {
		ret.W(idClause)
	}
	ret.W("\tout := util.ValueMap{}")
	grpcAddSection(ret, "list", 1)
	if ga.Grp == nil {
		ret.W("\tret, err := appState.Services.%s.List(p.Ctx, nil, &filter.Params{}%s, p.Logger)", m.Proper(), suffix)
	} else {
		ret.W("\tret, err := appState.Services.%s.Get%s(p.Ctx, nil, %s, &filter.Params{}%s, p.Logger)", m.Proper(), ga.APISuffix(), ga.Grp.Camel(), suffix)
	}
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\tout[%q] = ret", "results")
	ret.W("\tprovider.SetOutput(p.TX, out)")
	ret.W("\treturn p.TX, nil")
	ret.W("}")
	return ret
}

func grpcSearch(m *model.Model, grpcArgs string, grpcRet string, ga *FileArgs) *golang.Block {
	ret := golang.NewBlock("grpcSearch", "func")
	ret.W("func %sSearch%s(%s) %s {", m.Proper(), ga.APISuffix(), grpcArgs, grpcRet)
	idClause, suffix := idClauseFor(m)
	if idClause != "" {
		ret.W(idClause)
	}
	ret.W("\tq, _ := provider.GetString(p.R, p.TX, \"q\")")
	ret.W("\tif q == \"\" {")
	ret.W("\t\treturn nil, errors.New(\"must provide [q] in request data\")")
	ret.W("\t}")
	grpcAddSection(ret, "search", 1)
	if ga.Grp == nil {
		ret.W("\tret, err := appState.Services.%s.Search(p.Ctx, q, nil, nil%s, p.Logger)", m.Proper(), suffix)
	} else {
		ret.W("\tret, err := appState.Services.%s.Search%s(p.Ctx, %s, q, nil, nil%s, p.Logger)", m.Proper(), ga.APISuffix(), ga.Grp.Camel(), suffix)
	}
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\tprovider.SetOutput(p.TX, ret)")
	ret.W("\treturn p.TX, nil")
	ret.W("}")
	return ret
}

func grpcDetail(m *model.Model, grpcArgs string, grpcRet string, g *golang.File, ga *FileArgs) (*golang.Block, error) {
	ret := golang.NewBlock("grpcDetail", "func")
	ret.W("func %sDetail%s(%s) %s {", m.Proper(), ga.APISuffix(), grpcArgs, grpcRet)
	idClause, suffix := idClauseFor(m)
	if idClause != "" {
		ret.W(idClause)
	}
	pks := m.PKs()
	ret.W("\t%s, err := %sParamsFromRequest(p)", strings.Join(pks.CamelNames(), ", "), m.Camel())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	grpcAddSection(ret, "detail", 1)
	if m.IsRevision() {
		hc := m.HistoryColumn()
		if hc == nil {
			return nil, errors.New("invalid history column")
		}
		g.AddImport(helper.AppImport("app/" + m.Package))
		ret.W("\tvar ret *%s.%s", m.Package, m.Proper())
		ret.W("\t%s, err := provider.GetRequestInt(p.R, %q)", hc.Camel(), hc.Camel())
		ret.W("\tif err == nil && %s > 0 {", hc.Camel())
		xargs := strings.Join(pks.CamelNames(), ", ")
		ret.W("\t\tret, err = appState.Services.%s.Get%s(p.Ctx, nil, %s, %s, p.Logger)", m.Proper(), hc.Proper(), xargs, hc.Camel())
		ret.W("\t} else {")
		ret.W("\t\tret, err = appState.Services.%s.Get(p.Ctx, nil, %s%s, p.Logger)", m.Proper(), strings.Join(pks.CamelNames(), ", "), suffix)
		ret.W("\t}")
	} else {
		ret.W("\tret, err := appState.Services.%s.Get(p.Ctx, nil, %s%s, p.Logger)", m.Proper(), strings.Join(pks.CamelNames(), ", "), suffix)
	}
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ga.AddStaticCheck("ret", ret, m, ga.Grp, "retrieve")
	grpcAddSection(ret, "postdetail", 1)
	ret.W("\tprovider.SetOutput(p.TX, ret)")
	ret.W("\treturn p.TX, nil")
	ret.W("}")
	return ret, nil
}
