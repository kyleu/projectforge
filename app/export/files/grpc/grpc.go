package grpc

import (
	"fmt"
	"sort"

	"github.com/kyleu/projectforge/app/export/files/helper"
	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

func GRPC(m *model.Model, args *model.Args) (file.Files, error) {
	fileArgs, err := GetGRPCFileArgs(m, args)
	if err != nil {
		return nil, errors.Wrap(err, "invalid arguments")
	}

	var ret file.Files
	for _, fa := range fileArgs {
		g, err := grpcFile(m, args, fa)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		ret = append(ret, g)
	}
	return ret, nil
}

func GetGRPCFileArgs(m *model.Model, args *model.Args) ([]*GRPCFileArgs, error) {
	grpcPackage := args.Config.GetStringOpt("grpcPackage")
	if grpcPackage == "" {
		return nil, errors.New("must provide [grpcPackage] in the export config")
	}
	grpcClass := args.Config.GetStringOpt("grpcClass")
	if grpcClass == "" {
		return nil, errors.New("must provide [grpcClass] in the export config")
	}
	cPkg, _ := util.StringSplitLast(grpcClass, '.', true)
	grpcGroups, _ := args.Config.GetMap("grpcGroups")
	if len(grpcGroups) == 0 {
		grpcGroups = util.ValueMap{"*": "*"}
	}

	var ret []*GRPCFileArgs
	for k, grpIface := range grpcGroups {
		grp, ok := grpIface.(string)
		if !ok {
			return nil, errors.New("grpcGroups values must be strings")
		}
		if grp == "*" {
			ret = append(ret, &GRPCFileArgs{Class: grpcClass, Pkg: grpcPackage, CPkg: cPkg, API: k, Grp: nil})
		} else {
			g := m.Columns.Get(grp)
			if g == nil {
				return nil, errors.Errorf("grpcGroups references missing column [%s]", grp)
			}
			if !g.HasTag("grouped") {
				return nil, errors.Errorf("grpcGroups references non-grouped column [%s]", grp)
			}
			ret = append(ret, &GRPCFileArgs{Class: grpcClass, Pkg: grpcPackage, CPkg: cPkg, API: k, Grp: g})
		}
	}
	sort.Slice(ret, func(i, j int) bool {
		l, r := ret[i], ret[j]
		return l.API < r.API
	})
	return ret, nil
}

func grpcFile(m *model.Model, args *model.Args, ga *GRPCFileArgs) (*file.File, error) {
	fn := m.Package
	if ga.API != "*" {
		fn += "by" + ga.API
	}
	g := golang.NewFile(ga.Pkg, []string{"app", ga.Pkg}, fn)
	g.AddImport(helper.ImpErrors, helper.ImpFilter)
	g.AddImport(helper.AppImport("app/lib/" + ga.CPkg))
	if ga.Grp == nil {
		g.AddImport(helper.AppImport("app/" + m.Package))
	}

	grpcArgs := fmt.Sprintf("p *%s.Params", ga.CPkg)
	grpcRet := fmt.Sprintf("(*%sTransaction, error)", ga.Class)

	g.AddBlocks(
		grpcList(m, grpcArgs, grpcRet, ga),
		grpcSearch(m, grpcArgs, grpcRet, ga),
		grpcDetail(m, grpcArgs, grpcRet, ga),
		grpcCall("Create", m, grpcArgs, grpcRet, ga),
		grpcCall("Update", m, grpcArgs, grpcRet, ga),
		grpcCall("Save", m, grpcArgs, grpcRet, ga),
		grpcDelete(m, grpcArgs, grpcRet, ga),
	)
	if ga.Grp == nil {
		b, err := grpcParamsFromRequest(m, ga.CPkg)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(grpcFromRequest(m), b)
	}
	return g.Render()
}
