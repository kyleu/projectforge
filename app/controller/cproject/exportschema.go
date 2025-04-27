package cproject

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"path"
	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
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
		schCollection, err := schemasFor(ps.Context, prj, ps.Logger, "tmp/schema")
		ps.SetTitleAndData(fmt.Sprintf("[%s] JSON Schema", prj.Key), schCollection)
		page := &vexport.JSONSchema{Project: prj, Args: prj.ExportArgs, Collection: schCollection}
		return controller.Render(r, as, page, ps, "projects", prj.Key, "JSON Schema")
	})
}

type JSMap map[string]*jsonschema.Schema

func schemasFor(ctx context.Context, prj *project.Project, logger util.Logger, extraPaths ...string) (*jsonschema.Collection, error) {
	ret := jsonschema.NewCollection(prj.Package)
	for _, x := range prj.ExportArgs.Enums {
		_, err := schemaForEnum(prj, ret, x)
		if err != nil {
			return nil, err
		}
	}
	for _, x := range prj.ExportArgs.Models {
		_, err := schemaForModel(prj, ret, x)
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

func schemaForEnum(prj *project.Project, sch *jsonschema.Collection, x *enum.Enum) (*jsonschema.Schema, error) {
	ret := sch.NewSchema(x.Name)
	return ret, nil
}

func schemaForModel(prj *project.Project, sch *jsonschema.Collection, x *model.Model) (*jsonschema.Schema, error) {
	ret := sch.NewSchema(x.Name)
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
