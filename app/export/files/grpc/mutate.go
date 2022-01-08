package grpc

import (
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

func grpcCall(k string, m *model.Model, grpcArgs string, grpcRet string, ga *GRPCFileArgs) *golang.Block {
	ret := golang.NewBlock("grpc"+k, "func")
	ret.W("func %s%s%s(%s) %s {", m.Proper(), k, ga.APISuffix(), grpcArgs, grpcRet)
	ret.W("\tmodel, err := %sFromRequest(p.R)", m.Camel())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	grpcAddSection(ret, strings.ToLower(k), 1)
	ga.AddStaticCheck("model", ret, ga.Grp)
	ret.W("\terr = appState.Services.%s.%s(p.Ctx, nil, model)", m.Proper(), k)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\tprovider.SetOutput(p.TX, model)")
	ret.W("\treturn p.TX, nil")
	ret.W("}")
	return ret
}

func grpcDelete(m *model.Model, grpcArgs string, grpcRet string, ga *GRPCFileArgs) *golang.Block {
	pks := m.PKs()
	ret := golang.NewBlock("grpcDelete", "func")
	ret.W("func %sDelete%s(%s) %s {", m.Proper(), ga.APISuffix(), grpcArgs, grpcRet)
	ret.W("\t%s, err := %sParamsFromRequest(p.R)", strings.Join(pks.Names(), ", "), m.Camel())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	suffix := ""
	if m.IsSoftDelete() {
		suffix = ", true"
	}
	ret.W("\torig, err := appState.Services.%s.Get(p.Ctx, nil, %s%s)", m.Proper(), strings.Join(pks.Names(), ", "), suffix)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	grpcAddSection(ret, "delete", 1)
	ga.AddStaticCheck("orig", ret, ga.Grp)
	ret.W("\terr = appState.Services.%s.Delete(p.Ctx, nil, %s)", m.Proper(), strings.Join(pks.Names(), ", "))
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
	ret.W("func %sFromRequest(r *provider.NuevoRequest) (*%s, error) {", m.Package, m.ClassRef())
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
