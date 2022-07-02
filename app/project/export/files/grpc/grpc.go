package grpc

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/util"
)

const updateKey = "Update"

func GRPC(m *model.Model, args *model.Args, addHeader bool) (file.Files, error) {
	fileArgs, err := GetGRPCFileArgs(m, args)
	if err != nil {
		return nil, errors.Wrap(err, "invalid arguments")
	}

	var ret file.Files
	for _, fa := range fileArgs {
		g, err := grpcFile(m, args, fa, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		ret = append(ret, g)
	}
	return ret, nil
}

func GetGRPCFileArgs(m *model.Model, args *model.Args) ([]*FileArgs, error) {
	grpcPackage := args.Config.GetStringOpt("grpcPackage")
	if grpcPackage == "" {
		return nil, errors.New("must provide [grpcPackage] in the export config")
	}
	grpcClass := args.Config.GetStringOpt("grpcClass")
	if grpcClass == "" {
		return nil, errors.New("must provide [grpcClass] in the export config")
	}
	cPkg, _ := util.StringSplitLast(grpcClass, '.', true)
	grpcGroups, _ := args.Config.ParseMap("grpcGroups", true, true)
	if len(grpcGroups) == 0 {
		grpcGroups = util.ValueMap{"*": "*"}
	}

	var ret []*FileArgs
	for k, grpIface := range grpcGroups {
		grp, ok := grpIface.(string)
		if !ok {
			return nil, errors.New("grpcGroups values must be strings")
		}
		if grp == "*" {
			ret = append(ret, &FileArgs{Class: grpcClass, Pkg: grpcPackage, CPkg: cPkg, API: k, Grp: nil})
		} else {
			g := m.Columns.Get(grp)
			if g == nil {
				continue
				// return nil, errors.Errorf("grpcGroups references missing column [%s]", grp)
			}
			if !g.HasTag("grouped") {
				return nil, errors.Errorf("grpcGroups references non-grouped column [%s]", grp)
			}
			ret = append(ret, &FileArgs{Class: grpcClass, Pkg: grpcPackage, CPkg: cPkg, API: k, Grp: g})
		}
	}
	slices.SortFunc(ret, func(l *FileArgs, r *FileArgs) bool {
		return l.API < r.API
	})
	return ret, nil
}

func grpcFile(m *model.Model, args *model.Args, ga *FileArgs, addHeader bool) (*file.File, error) {
	fn := m.Package
	if ga.API != "*" {
		fn += "by" + ga.API
	}
	g := golang.NewFile(ga.Pkg, []string{"app", ga.Pkg}, fn)
	g.AddImport(helper.ImpErrors, helper.ImpFilter, helper.ImpAppUtil)
	g.AddImport(helper.AppImport("app/lib/" + ga.CPkg))
	if ga.Grp == nil {
		g.AddImport(helper.AppImport("app/" + m.PackageWithGroup("")))
	}

	grpcArgs := fmt.Sprintf("p *%s.Params", ga.CPkg)
	grpcRet := fmt.Sprintf("(*%sTransaction, error)", ga.Class)

	g.AddBlocks(grpcList(m, grpcArgs, grpcRet, ga))
	if len(m.Search) > 0 {
		g.AddBlocks(grpcSearch(m, grpcArgs, grpcRet, ga))
	}
	if detail, err := grpcDetail(m, grpcArgs, grpcRet, g, ga); err == nil {
		g.AddBlocks(detail)
	} else {
		return nil, err
	}
	g.AddBlocks(
		grpcCall("Create", m, false, grpcArgs, grpcRet, ga),
		grpcCall(updateKey, m, true, grpcArgs, grpcRet, ga),
		grpcCall("Save", m, true, grpcArgs, grpcRet, ga),
		grpcDelete(m, grpcArgs, grpcRet, ga),
	)
	if ga.Grp == nil {
		b, err := grpcParamsFromRequest(m, grpcArgs, g)
		if err != nil {
			return nil, err
		}
		g.AddBlocks(grpcFromRequest(m, grpcArgs), b)
	}
	return g.Render(addHeader)
}
