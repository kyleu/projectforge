package metaschema

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

func ExportModel(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*model.Model, error) {
	if sch == nil {
		return nil, errors.New("nil schema provided for model")
	}
	_, err := exportGetType(sch, "object")
	if err != nil {
		return nil, err
	}
	n, pkg, grp := parseID(sch.ID)
	ret := &model.Model{Name: n, Package: pkg, Group: grp, Description: sch.Description}

	for _, propKey := range sch.Properties.Order {
		col, e := ExportColumn(propKey, sch, coll, args)
		if e != nil {
			return nil, e
		}
		ret.Columns = append(ret.Columns, col)
	}

	if x := sch.Metadata.GetMapOpt("config"); len(x) > 0 {
		ret.Config = x
	}
	if x := sch.Metadata.GetStringOpt("icon"); x != "" {
		ret.Icon = x
	}
	if x := sch.Metadata["ordering"]; x != nil {
		var ords filter.Orderings
		if e := util.CycleJSON(x, &ords); e != nil {
			return nil, e
		}
		ret.Ordering = ords
	}
	if x := sch.Metadata.GetStringOpt("plural"); x != "" {
		ret.PluralOverride = x
	}
	if x := sch.Metadata.GetStringOpt("route"); len(x) > 0 {
		ret.RouteOverride = x
	}
	if x := sch.Metadata.GetStringArrayOpt("search"); len(x) > 0 {
		ret.Search = x
	}
	if x := sch.Metadata.GetIntOpt("sortIndex"); x != 0 {
		ret.SortIndex = x
	}
	if x := sch.Metadata.GetStringOpt("table"); len(x) > 0 {
		ret.TableOverride = x
	}
	if x := sch.Metadata.GetStringArrayOpt("tags"); len(x) > 0 {
		ret.Tags = x
	}
	if x := sch.Metadata.GetStringOpt("title"); x != "" {
		ret.TitleOverride = x
	}
	if x := sch.Metadata.GetStringOpt("view"); x != "" {
		ret.View = x
	}
	return ret, nil
}
