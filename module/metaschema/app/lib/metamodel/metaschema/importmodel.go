package metaschema

import (
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/lib/jsonschema"
	"{{{ .Package }}}/app/lib/metamodel"
	"{{{ .Package }}}/app/lib/metamodel/model"
	"{{{ .Package }}}/app/util"
)

func ImportModel(sch *jsonschema.Schema, coll *jsonschema.Collection, args *metamodel.Args) (*model.Model, error) {
	if sch == nil {
		return nil, errors.New("nil schema provided for model")
	}
	_, err := exportGetType(sch, "object")
	if err != nil {
		return nil, err
	}
	n, pkg, grp := parseID(sch.ID())
	ret := &model.Model{Name: n, Package: pkg, Group: grp, Description: sch.Description}

	for _, propKey := range sch.Properties.Order {
		col, e := ImportColumn(propKey, sch, coll, args)
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
	if x := sch.Metadata.GetStringOpt("route"); x != "" {
		ret.RouteOverride = x
	}
	if x := sch.Metadata.GetStringOpt("schema"); x != "" {
		ret.Schema = x
	}
	if x := sch.Metadata.GetStringArrayOpt("search"); len(x) > 0 {
		ret.Search = x
	}
	if x := sch.Metadata.GetIntOpt("sortIndex"); x != 0 {
		ret.SortIndex = x
	}
	if x := sch.Metadata.GetStringOpt("table"); x != "" {
		ret.TableOverride = x
	}
	if x := sch.Metadata.GetStringArrayOpt("tags"); len(x) > 0 {
		ret.Tags = x
	}
	if sch.Title != "" && sch.Title != n {
		ret.TitleOverride = sch.Title
	}
	if x := sch.Metadata.GetStringOpt("view"); x != "" {
		ret.View = x
	}
	if len(sch.Examples) > 0 {
		ret.SeedData = util.ArrayFromAny[[]any](sch.Examples)
	}
	if e := parseExtra(sch.Metadata, ret); e != nil {
		return nil, e
	}
	return ret, nil
}

func parseExtra(md util.ValueMap, ret *model.Model) error {
	if x, ok := md["relations"]; ok {
		var rels model.Relations
		if err := util.CycleJSON(x, &rels); err != nil {
			return err
		}
		ret.Relations = rels
	}
	if x, ok := md["indexes"]; ok {
		var idxs model.Indexes
		if err := util.CycleJSON(x, &idxs); err != nil {
			return err
		}
		ret.Indexes = idxs
	}
	if x, ok := md["links"]; ok {
		var links model.Links
		if err := util.CycleJSON(x, &links); err != nil {
			return err
		}
		ret.Links = links
	}
	if x, ok := md["imports"]; ok {
		var imps model.Imports
		if err := util.CycleJSON(x, &imps); err != nil {
			return err
		}
		ret.Imports = imps
	}
	return nil
}
