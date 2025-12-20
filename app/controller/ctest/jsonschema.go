package ctest

import (
	"net/http"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel/metaschema"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views"
	"projectforge.dev/projectforge/views/layout"
	"projectforge.dev/projectforge/views/vjsonschema"
)

const schemaPath = "tmp/schemastore/src"

func schemaInit(ps *cutil.PageState) (layout.Page, error) {
	fs, err := filesystem.NewFileSystem(schemaPath, true, "")
	if err != nil {
		return nil, err
	}
	ret := util.NewStringSlice()
	fin := func(msg string, args ...any) (layout.Page, error) {
		ret.Pushf(msg, args...)
		ps.SetTitleAndData("JSON Schema Init", ret.Slice)
		return &views.Debug{}, nil
	}

	if fs.IsDir(".") {
		return fin("We're all good!")
	}

	if !fs.IsDir(".") {
		ret.Pushf("cloning [schemastore] repo...")
		_, _, err := util.RunProcessSimple("git clone https://github.com/SchemaStore/schemastore.git", "./tmp")
		if err != nil {
			return fin("Error cloning repo: %+v", err)
		}
		ret.Pushf("completed clone of [schemastore] repo")
	}
	return fin("Initialization complete")
}

func schemaTest(ps *cutil.PageState) (layout.Page, error) {
	fs, err := filesystem.NewFileSystem(schemaPath, true, "")
	if err != nil {
		return nil, err
	}
	filenames := util.ArraySorted(fs.ListJSON("schemas/json", nil, true, ps.Logger))
	var ret metaschema.SchemaTestFiles
	ret, _ = util.AsyncCollect(filenames, func(f string) (*metaschema.SchemaTestFile, error) {
		return metaschema.LoadSchemaTestFile(f, fs, ps.Logger), nil
	}, ps.Logger)
	ret.Sort()
	ps.SetTitleAndData("JSON Schema Test", ret)
	page := &vjsonschema.SchemaList{Schemata: ret, ShowContent: false}
	return page, nil
}

func JSONSchemaTestFile(w http.ResponseWriter, r *http.Request) {
	controller.Act("json.schema.test.file", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		filename, err := cutil.PathString(r, "filename", false)
		if err != nil {
			return "", err
		}
		fs, err := filesystem.NewFileSystem(schemaPath, true, "")
		if err != nil {
			return "", err
		}
		ret := metaschema.LoadSchemaTestFile(filename, fs, ps.Logger)
		ps.SetTitleAndData("["+filename+"] JSON Schema Test", ret)
		return controller.Render(r, as, &vjsonschema.SchemaDetail{Schema: ret}, ps, "Tests||/test", "JSON Schema||/test/jsonschema-test", filename)
	})
}
