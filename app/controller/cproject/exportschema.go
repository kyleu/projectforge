package cproject

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"slices"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel/metaschema"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportJSONSchema(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.json.schema", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		schCollection, err := schemasFor(ps.Context, prj, []string{"tmp/schema"}, ps.Logger)
		ps.SetTitleAndData(fmt.Sprintf("[%s] JSON Schema", prj.Key), schCollection)
		page := &vexport.JSONSchemaCollection{Project: prj, Args: prj.ExportArgs, Collection: schCollection}
		return controller.Render(r, as, page, ps, "projects", prj.Key, "JSON Schema")
	})
}

func ProjectExportModelJSONSchema(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.json.model.schema", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, x, err := exportLoadModel(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		schCollection, err := schemasFor(ps.Context, prj, nil, ps.Logger, x.Name)
		ps.SetTitleAndData(fmt.Sprintf("[%s] JSON Schema", prj.Key), schCollection)
		page := &vexport.JSONSchemaModel{Project: prj, Model: x, Collection: schCollection}
		return controller.Render(r, as, page, ps, "projects", prj.Key, x.Title()+" JSON Schema")
	})
}

func ProjectExportEnumJSONSchema(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.json.enum.schema", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, x, err := exportLoadEnum(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		schCollection, err := schemasFor(ps.Context, prj, nil, ps.Logger, x.Name)
		ps.SetTitleAndData(fmt.Sprintf("[%s] JSON Schema", prj.Key), schCollection)
		page := &vexport.JSONSchemaEnum{Project: prj, Enum: x, Collection: schCollection}
		return controller.Render(r, as, page, ps, "projects", prj.Key, x.Title()+" JSON Schema")
	})
}

func schemasFor(ctx context.Context, prj *project.Project, extraPaths []string, logger util.Logger, filter ...string) (*jsonschema.Collection, error) {
	ret := jsonschema.NewCollection(prj.Package)
	for _, x := range prj.ExportArgs.Enums {
		if len(filter) > 0 && !slices.Contains(filter, x.Name) {
			continue
		}
		_, err := metaschema.EnumSchema(x, ret, prj.ExportArgs)
		if err != nil {
			return nil, err
		}
	}
	for _, x := range prj.ExportArgs.Models {
		if len(filter) > 0 && !slices.Contains(filter, x.Name) {
			continue
		}
		_, err := metaschema.ModelSchema(x, ret, prj.ExportArgs)
		if err != nil {
			return nil, err
		}
	}
	for _, x := range extraPaths {
		err := parseExtraPath(ctx, x, ret, logger)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse [%s]", x)
		}

	}
	return ret, nil
}

func parseExtraPath(ctx context.Context, pth string, coll *jsonschema.Collection, logger util.Logger) error {
	fs, _ := filesystem.NewFileSystem(".", true, "")
	if fs.IsDir(pth) {
		files := fs.ListJSON(pth, nil, false, logger)
		for _, fn := range files {
			if err := parseExtraPath(ctx, path.Join(pth, fn), coll, logger); err != nil {
				return err
			}
		}
	} else {
		b, err := fs.ReadFile(pth)
		if err != nil {
			return err
		}
		sch, err := jsonschema.FromJSON(b)
		if err != nil {
			return err
		}
		coll.AddSchema(sch)
	}
	return nil
}
