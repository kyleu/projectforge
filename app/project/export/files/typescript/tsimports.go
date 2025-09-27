package typescript

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/util"
)

func tsModelImportsColumn(col *model.Column, args *metamodel.Args, str model.StringProvider) ([]string, error) {
	var ret TSImports
	if e, _ := model.AsEnum(col.Type); e != nil {
		if en := args.Enums.Get(e.Ref); en != nil {
			pth := str.RelativePath(en.GroupAndPackage(), en.Camel())
			op := fmt.Sprintf(`%s%s`, util.Choose(col.Nullable, "parse", "get"), en.Proper())
			ret = ret.With(newImport(en.Proper(), pth), newImport(op, pth))
		}
	}
	r, rm, _ := helper.LoadRef(col, args.Models, args.Events, args.ExtraTypes)
	if rm == nil {
		if col.Metadata != nil {
			if tsImport := col.Metadata.GetStringOpt("tsImport"); tsImport != "" {
				ret = ret.With(newImport(r.K, tsImport))
			}
		}
	} else {
		if tsImport := rm.ConfigMap().GetStringOpt("tsImport"); tsImport != "" {
			ret = ret.With(newImport(r.K, tsImport))
		} else if rm.PackageWithGroup("") != str.PackageWithGroup("") {
			ret = ret.With(newImport(r.K, str.RelativePath(rm.GroupAndPackage(), rm.Camel())))
		} else {
			ret = ret.With(newImport(r.K, fmt.Sprintf("./%s", rm.Camel())))
		}
	}
	if col.Type.Key() == types.KeyNumeric {
		pth := fmt.Sprintf("%snumeric/numeric", util.StringRepeat("../", str.GroupLen()+1))
		ret = ret.With(newImport("Numeric", pth), newImport("NumericSource", pth))
	}
	return ret.Strings(), nil
}

func tsModelImports(args *metamodel.Args, cols model.Columns, str model.StringProvider) ([]string, error) {
	ret := &util.StringSlice{}
	add := func(s string, args ...any) {
		ret.PushUnique(fmt.Sprintf(s, args...))
	}
	relPath := util.StringRepeat("../", str.GroupLen()+1)
	add(`import { Parse } from "%sparse";`, relPath)
	for _, col := range cols {
		x, err := tsModelImportsColumn(col, args, str)
		if err != nil {
			return nil, err
		}
		ret.PushUnique(x...)
	}
	return ret.Slice, nil
}
