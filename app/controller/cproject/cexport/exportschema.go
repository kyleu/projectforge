package cexport

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cproject"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel/jsonload"
	"projectforge.dev/projectforge/app/lib/metamodel/metaschema"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vjsonschema"
)

func ProjectExportJSONSchema(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.json.schema", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := cproject.GetProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		schCollection, err := metaschema.LoadSchemas(ps.Context, prj.Key, prj.ExportArgs, nil, ps.Logger)
		if err != nil {
			return "", err
		}
		oldResults, err := metaschema.ImportArgs(schCollection, prj.ExportArgs)
		if err != nil {
			return "", err
		}
		newValidation := &jsonload.Validation{Collection: schCollection}
		_, newResults, err := newValidation.Export(ps.Context, ps.Logger)
		if err != nil {
			return "", err
		}
		unrelated := lo.Reject(schCollection.Schemas(), func(x *jsonschema.Schema, _ int) bool {
			return strings.Contains(x.Comment, util.AppName)
		})
		ps.SetTitleAndData(fmt.Sprintf("[%s] JSON Schema", prj.Key), schCollection)
		page := &vjsonschema.CollectionDetail{
			BaseURL: prj.WebPath(), Args: prj.ExportArgs, Collection: schCollection, OldResults: oldResults, NewResults: newResults, Unrelated: unrelated,
		}
		return controller.Render(r, as, page, ps, "projects", prj.Key, "JSON Schema")
	})
}

func ProjectExportWriteJSONSchema(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.write.json.schema", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := cproject.GetProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		schCollection, err := metaschema.LoadSchemas(ps.Context, prj.Key, prj.ExportArgs, nil, ps.Logger)
		if err != nil {
			return "", err
		}
		fs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		if !fs.IsDir("tmp/schema") {
			err = fs.CreateDirectory("tmp/schema")
			if err != nil {
				return "", err
			}
		}
		for _, sch := range schCollection.SchemaMap {
			id := sch.ID()
			if strings.Contains(id, "/") {
				_, id = util.StringCutLast(id, '/', true)
			}
			fn := "tmp/schema/" + id
			err = fs.WriteFile(fn, util.ToJSONBytes(sch, true), filesystem.DefaultMode, true)
			if err != nil {
				return "", err
			}
		}
		return controller.FlashAndRedir(true, "wrote JSON Schema files", prj.WebPath()+"/export/jsonschema", ps)
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
		tgt, err := metaschema.ImportModel(sch, schCollection, prj.ExportArgs)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] JSON Schema", prj.Key), schCollection)
		page := &vjsonschema.ModelDetail{BaseURL: prj.WebPath() + "/export/models", Model: x, Collection: schCollection, Result: tgt}
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
		tgt, err := metaschema.ImportEnum(sch, schCollection, prj.ExportArgs)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] JSON Schema", prj.Key), schCollection)
		page := &vjsonschema.EnumDetail{BaseURL: prj.WebPath() + "/export/enums", Enum: x, Collection: schCollection, Result: tgt}
		return controller.Render(r, as, page, ps, "projects", prj.Key, x.Title()+" JSON Schema")
	})
}
