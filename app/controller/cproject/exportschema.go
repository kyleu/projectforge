package cproject

import (
	"fmt"
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel/metaschema"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportJSONSchema(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.json.schema", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := getProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		schCollection, err := metaschema.LoadSchemas(ps.Context, prj.Key, prj.ExportArgs, []string{"tmp/schema"}, ps.Logger)
		if err != nil {
			return "", err
		}
		results, err := metaschema.ExportArgs(schCollection, prj.ExportArgs)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] JSON Schema", prj.Key), schCollection)
		page := &vexport.JSONSchemaCollection{Project: prj, Args: prj.ExportArgs, Collection: schCollection, Results: results}
		return controller.Render(r, as, page, ps, "projects", prj.Key, "JSON Schema")
	})
}

func ProjectExportModelJSONSchema(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.json.model.schema", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, x, err := exportLoadModel(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		schCollection, err := metaschema.LoadSchemas(ps.Context, prj.Key, prj.ExportArgs, nil, ps.Logger, x.Name)
		if err != nil {
			return "", err
		}
		sch := schCollection.GetSchema(x.ID())
		tgt, err := metaschema.ExportModel(sch, schCollection, prj.ExportArgs)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] JSON Schema", prj.Key), schCollection)
		page := &vexport.JSONSchemaModel{Project: prj, Model: x, Collection: schCollection, Result: tgt}
		return controller.Render(r, as, page, ps, "projects", prj.Key, x.Title()+" JSON Schema")
	})
}

func ProjectExportEnumJSONSchema(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.json.enum.schema", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, x, err := exportLoadEnum(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		schCollection, err := metaschema.LoadSchemas(ps.Context, prj.Key, prj.ExportArgs, nil, ps.Logger, x.Name)
		if err != nil {
			return "", err
		}
		sch := schCollection.GetSchema(x.ID())
		tgt, err := metaschema.ExportEnum(sch, schCollection, prj.ExportArgs)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] JSON Schema", prj.Key), schCollection)
		page := &vexport.JSONSchemaEnum{Project: prj, Enum: x, Collection: schCollection, Result: tgt}
		return controller.Render(r, as, page, ps, "projects", prj.Key, x.Title()+" JSON Schema")
	})
}
