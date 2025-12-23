package jsonload

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

func exportObject(ctx context.Context, sch *jsonschema.Schema, ret *metamodel.Args, s *jsonschema.Collection, logs *util.StringSlice) error {
	logs.Pushf("exporting object schema [%s]", sch.ID())
	name, pkg, grp, err := extractPath(sch)
	if err != nil {
		return errors.Wrapf(err, "unable to process [path] from [%s]", sch.ID())
	}
	md := sch.GetMetadata()

	ords, err := fromMD[filter.Orderings](sch.Unknown, "ordering")
	if err != nil {
		return errors.Wrapf(err, "unable to process [ordering] for [%s]", sch.ID())
	}
	cols, err := parseProperties(sch)
	if err != nil {
		return errors.Wrapf(err, "unable to process [columns] for [%s]", sch.ID())
	}
	rels, err := fromMD[model.Relations](sch.Unknown, "relations")
	if err != nil {
		return errors.Wrapf(err, "unable to process [relations] for [%s]", sch.ID())
	}
	idxs, err := fromMD[model.Indexes](sch.Unknown, "indexes")
	if err != nil {
		return errors.Wrapf(err, "unable to process [indexes] for [%s]", sch.ID())
	}
	seed, err := fromMD[[][]any](sch.Unknown, "seed")
	if err != nil {
		return errors.Wrapf(err, "unable to process [seed] for [%s]", sch.ID())
	}
	links, err := fromMD[model.Links](sch.Unknown, "links")
	if err != nil {
		return errors.Wrapf(err, "unable to process [links] for [%s]", sch.ID())
	}
	imports, err := fromMD[model.Imports](sch.Unknown, "imports")
	if err != nil {
		return errors.Wrapf(err, "unable to process [imports] for [%s]", sch.ID())
	}

	m := &model.Model{
		Name: name, Package: pkg, Group: grp,

		Schema:         md.GetStringOpt("schema"),
		Description:    sch.Description,
		Icon:           md.GetStringOpt("icon"),
		Ordering:       ords,
		SortIndex:      md.GetIntOpt("sortIndex"),
		View:           md.GetStringOpt("view"),
		Search:         md.GetStringArrayOpt("search"),
		Tags:           md.GetStringArrayOpt("tags"),
		TitleOverride:  md.GetStringOpt("title"),
		PluralOverride: md.GetStringOpt("plural"),
		ProperOverride: md.GetStringOpt("proper"),
		TableOverride:  md.GetStringOpt("table"),
		RouteOverride:  md.GetStringOpt("route"),
		Config:         md.GetMapOpt("config"),
		Columns:        cols,
		Relations:      rels,
		Indexes:        idxs,
		SeedData:       seed,
		Links:          links,
		Imports:        imports,
	}
	ret.Models = append(ret.Models, m)
	return nil
}
