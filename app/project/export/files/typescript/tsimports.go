package typescript

import (
	"fmt"

	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/files/helper"
	"projectforge.dev/projectforge/app/util"
)

const localRef = "./"

func relativePath(sp metamodel.StringProvider, rGroup []string, extra ...string) string {
	mGroup := sp.GroupAndPackage()
	commonPrefix := 0
	for i := 0; i < len(mGroup) && i < len(rGroup) && mGroup[i] == rGroup[i]; i++ {
		commonPrefix++
	}
	upLevels := len(mGroup) - commonPrefix
	var pathParts []string
	for i := commonPrefix; i < len(rGroup); i++ {
		pathParts = append(pathParts, rGroup[i])
	}
	pathParts = append(pathParts, extra...)
	p := util.StringRepeat("../", upLevels)
	if upLevels == 0 {
		p = localRef
	}
	return fmt.Sprintf("%s%s", p, util.StringJoin(pathParts, "/"))
}

func tsModelImportsColumn(col *model.Column, args *metamodel.Args, str metamodel.StringProvider) ([]string, error) {
	var ret TSImports
	if e, _ := model.AsEnum(col.Type); e != nil {
		if en := args.Enums.Get(e.Ref); en != nil {
			pth := relativePath(str, en.GroupAndPackage(), en.Kebab())
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
			ret = ret.With(newImport(r.K, relativePath(str, rm.GroupAndPackage(), rm.Kebab())))
		} else {
			ret = ret.With(newImport(r.K, fmt.Sprintf("./%s", rm.Kebab())))
		}
	}
	if col.Type.Key() == types.KeyNumeric {
		pth := "@numeric/numeric"
		ret = ret.With(newImport("Numeric", pth), newImport("NumericSource", pth))
	}
	if col.Type.Key() == types.KeyNumericMap {
		pth := "@numeric/numeric-map"
		ret = ret.With(newImport("NumericMap", pth), newImport("parseNumericMap", pth))
	}
	return ret.Strings(), nil
}

func tsModelImports(args *metamodel.Args, cols model.Columns, str metamodel.StringProvider) ([]string, error) {
	ret := &util.StringSlice{}
	add := func(s string, args ...any) {
		ret.PushUnique(fmt.Sprintf(s, args...))
	}
	relPath := util.StringRepeat("../", str.GroupLen()+1)
	if str.GroupLen() == 0 {
		relPath = localRef
	}
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
