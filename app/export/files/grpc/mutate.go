package grpc

import (
	"strings"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

func grpcCall(k string, m *model.Model, isUpdate bool, grpcArgs string, grpcRet string, ga *FileArgs) *golang.Block {
	ret := golang.NewBlock("grpc"+k, "func")
	ret.W("func %s%s%s(%s) %s {", m.Proper(), k, ga.APISuffix(), grpcArgs, grpcRet)
	ret.W("\tmodel, err := %sFromRequest(p.R, %t)", m.Camel(), isUpdate)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	grpcAddSection(ret, strings.ToLower(k), 1)
	ga.AddStaticCheck("model", ret, m, ga.Grp, strings.ToLower(k))
	ret.W("\terr = appState.Services.%s.%s(p.Ctx, nil, model)", m.Proper(), k)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	if k == "Create" {
		grpcAddSection(ret, "post"+strings.ToLower(k), 1)
	}
	ret.W("\tprovider.SetOutput(p.TX, model)")
	ret.W("\treturn p.TX, nil")
	ret.W("}")
	return ret
}

func grpcDelete(m *model.Model, grpcArgs string, grpcRet string, ga *FileArgs) *golang.Block {
	pks := m.PKs()
	ret := golang.NewBlock("grpcDelete", "func")
	ret.W("func %sDelete%s(%s) %s {", m.Proper(), ga.APISuffix(), grpcArgs, grpcRet)
	ret.W("\t%s, err := %sParamsFromRequest(p.R)", strings.Join(pks.CamelNames(), ", "), m.Camel())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	ret.W("\torig, err := appState.Services.%s.Get(p.Ctx, nil, %s%s)", m.Proper(), strings.Join(pks.CamelNames(), ", "), m.SoftDeleteSuffix())
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, err")
	ret.W("\t}")
	grpcAddSection(ret, "delete", 1)
	ga.AddStaticCheck("orig", ret, m, ga.Grp, "delete")
	ret.W("\terr = appState.Services.%s.Delete(p.Ctx, nil, %s)", m.Proper(), strings.Join(pks.CamelNames(), ", "))
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
	ret.W("func %sFromRequest(r *provider.NuevoRequest, isUpdate bool) (*%s, error) {", m.Camel(), m.ClassRef())
	ret.W("\tinput := provider.GetRequest(r, \"%s\")", m.Camel())
	ret.W("\tif input == nil {")
	ret.W("\t\treturn nil, errors.New(\"must provide [%s] in request data\")", m.Camel())
	ret.W("\t}")

	ret.W("\tm, ok := input.(map[string]interface{})")
	ret.W("\tif !ok {")
	ret.W("\t\treturn nil, errors.New(\"field [%s] must be an object\")", m.Camel())
	ret.W("\t}")

	ret.W("\tret, err := %s.FromMap(m, true)", m.Package)
	ret.W("\tif err != nil {")
	ret.W("\t\treturn nil, errors.Wrap(err, \"unable to parse input data as %s\")", m.Camel())
	ret.W("\t}")
	ret.W("\t// $PF_SECTION_START(validate)$")
	ret.W("\t// $PF_SECTION_END(validate)$")
	ret.W("\treturn ret, nil")
	ret.W("}")
	return ret
}
