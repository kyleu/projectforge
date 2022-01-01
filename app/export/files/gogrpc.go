package files

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func GRPC(m *model.Model, args *model.Args) (*file.File, error) {
	grpcPackage := args.Config.GetStringOpt("grpcPackage")
	if grpcPackage == "" {
		return nil, errors.New("must provide [grpcPackage] in the export config")
	}
	grpcClass := args.Config.GetStringOpt("grpcClass")
	if grpcClass == "" {
		return nil, errors.New("must provide [grpcClass] in the export config")
	}
	cPkg, _ := util.StringSplitLast(grpcClass, '.', true)

	g := golang.NewFile(grpcPackage, []string{"app", grpcPackage}, m.Package)
	g.AddImport(golang.ImportTypeExternal, "github.com/pkg/errors")
	g.AddImport(golang.ImportTypeExternal, "go.uber.org/zap")
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/lib/"+cPkg)
	g.AddImport(golang.ImportTypeApp, "{{{ .Package }}}/app/"+m.Package)

	grpcArgs := fmt.Sprintf("p *provider.Params")
	grpcRet := fmt.Sprintf("(*%sTransaction, error)", grpcClass)

	g.AddBlocks(
		grpcList(m, grpcArgs, grpcRet), grpcDetail(m, grpcArgs, grpcRet), grpcCall("Add", m, grpcArgs, grpcRet),
		grpcCall("Update", m, grpcArgs, grpcRet), grpcCall("Save", m, grpcArgs, grpcRet), grpcDelete(m, grpcArgs, grpcRet),
		grpcFromRequest(m), grpcParamsFromRequest(m, cPkg),
	)
	return g.Render(), nil
}

func grpcList(m *model.Model, grpcArgs string, grpcRet string) *golang.Block {
	ret := golang.NewBlock("grpcList", "func")
	ret.W("func %sList(%s) %s {", m.Proper(), grpcArgs, grpcRet)
	ret.W("\tret, err := appState.Services.%s.List(p.Ctx, nil, nil)", m.Proper())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\tprovider.SetOutput(p.TX, ret)")
	ret.W("\treturn p.TX, nil")
	ret.W("}")
	return ret
}

func grpcDetail(m *model.Model, grpcArgs string, grpcRet string) *golang.Block {
	ret := golang.NewBlock("grpcDetail", "func")
	ret.W("func %sDetail(%s) %s {", m.Proper(), grpcArgs, grpcRet)
	ret.W("\tkey, revision, err := %sParamsFromRequest(p.R, p.Logger)", m.Camel())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\tret, err := appState.Services.%s.Get(p.Ctx, nil, key, revision)", m.Proper())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\tprovider.SetOutput(p.TX, ret)")
	ret.W("\treturn p.TX, nil")
	ret.W("}")
	return ret
}

func grpcCall(k string, m *model.Model, grpcArgs string, grpcRet string) *golang.Block {
	ret := golang.NewBlock("grpcAdd", "func")
	ret.W("func %s%s(%s) %s {", m.Proper(), k, grpcArgs, grpcRet)
	ret.W("\tmodel, err := %sFromRequest(p.R, p.Logger)", m.Camel())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\terr = appState.Services.%s.%s(p.Ctx, nil, model)", m.Proper(), k)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\tprovider.SetOutput(p.TX, model)")
	ret.W("\treturn p.TX, nil")
	ret.W("}")
	return ret
}

func grpcDelete(m *model.Model, grpcArgs string, grpcRet string) *golang.Block {
	ret := golang.NewBlock("grpcDelete", "func")
	ret.W("func %sDelete(%s) %s {", m.Proper(), grpcArgs, grpcRet)
	ret.W("\tkey, revision, err := %sParamsFromRequest(p.R, p.Logger)", m.Camel())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\torig, err := appState.Services.%s.Get(p.Ctx, nil, key, revision)", m.Proper())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\terr = appState.Services.%s.Delete(p.Ctx, nil, key, revision)", m.Proper())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\tprovider.SetOutput(p.TX, orig)")
	ret.W("\treturn p.TX, nil")
	ret.W("}")
	return ret
}

func grpcFromRequest(m *model.Model) *golang.Block {
	ret := golang.NewBlock("grpcFromRequest", "func")
	ret.W("func %sFromRequest(r *provider.NuevoRequest, logger *zap.SugaredLogger) (*%s, error) {", m.Package, m.ClassRef())
	ret.W("\tinput := provider.GetRequest(r, \"%s\")", m.Camel())
	ret.W("\tif input == nil {")
	ret.W("\t\treturn nil, errors.New(\"must provide [%s] in request data\")", m.Camel())
	ret.W("\t}")

	ret.W("\tm, ok := input.(map[string]interface{})")
	ret.W("\tif !ok {")
	ret.W("\t\treturn nil, errors.New(\"field [%s] must be an object\")", m.Camel())
	ret.W("\t}")

	ret.W("\tret, err := %s.FromMap(m, true)", m.Camel())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrap(err, \"unable to parse input data as %s\")", m.Camel())
	ret.W("\t}")
	ret.W("\treturn ret, nil")
	ret.W("}")
	return ret
}

func grpcParamsFromRequest(m *model.Model, cPkg string) *golang.Block {
	ret := golang.NewBlock("grpcParamsFromRequest", "func")
	ret.W("func %sParamsFromRequest(r *provider.NuevoRequest, logger *zap.SugaredLogger) (string, int, error) {", m.Camel())
	pks := m.Columns.PKs()
	zeroVals := strings.Join(pks.ZeroVals(), ", ")
	for _, col := range pks {
		grpcArgFor(col, ret, zeroVals)
	}
	ret.W("\treturn %s, nil", strings.Join(pks.CamelNames(), ", "))
	ret.W("}")
	return ret
}
