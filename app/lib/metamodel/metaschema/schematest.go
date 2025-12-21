package metaschema

import (
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/jsonschema"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/util"
)

type SchemaTestFile struct {
	Filename   string                 `json:"filename,omitempty"`
	Content    string                 `json:"content,omitempty"`
	Schema     *jsonschema.Schema     `json:"schema,omitempty"`
	Collection *jsonschema.Collection `json:"collection,omitempty"`
	Logs       []string               `json:"logs,omitempty"`
	Error      string                 `json:"error,omitempty"`
	Args       *metamodel.Args        `json:"args,omitempty"`
	ArgsError  string                 `json:"argsError,omitempty"`
}

func (s *SchemaTestFile) WithError(err error, showStack bool) *SchemaTestFile {
	if showStack {
		s.Error = fmt.Sprintf("%+v", err)
	} else {
		s.Error = err.Error()
	}
	return s
}

func (s *SchemaTestFile) Size() string {
	return util.ByteSizeSI(int64(len(s.Content)))
}

type SchemaTestFiles []*SchemaTestFile

func (s SchemaTestFiles) Sort() {
	slices.SortFunc(s, func(a, b *SchemaTestFile) int {
		return strings.Compare(strings.ToLower(a.Filename), strings.ToLower(b.Filename))
	})
}

func (s SchemaTestFiles) WithError() SchemaTestFiles {
	return lo.Filter(s, func(f *SchemaTestFile, _ int) bool {
		return f.Error != ""
	})
}

func (s SchemaTestFiles) WithoutError() SchemaTestFiles {
	return lo.Filter(s, func(f *SchemaTestFile, _ int) bool {
		return f.Error == ""
	})
}

func LoadSchemaTestFile(filename string, fs filesystem.FileLoader, logger util.Logger) *SchemaTestFile {
	ret := &SchemaTestFile{Filename: filename}
	t := util.TimerStart()
	json, err := fs.ReadFile(fmt.Sprintf("schemas/json/%s.json", filename))
	if err != nil {
		return ret.WithError(err, true)
	}
	ret.Content = string(json)

	sch := &jsonschema.Schema{}
	if err = util.FromJSON(json, &sch); err != nil {
		return ret.WithError(err, false)
	}
	ret.Schema = sch

	coll := jsonschema.NewCollection()
	if err := coll.AddSchemaExpanded(sch); err != nil {
		return ret.WithError(err, false)
	}
	ret.Collection = coll

	args := &metamodel.Args{}
	args, err = ImportArgs(coll, args)
	if err != nil {
		ret.ArgsError = err.Error()
	}
	ret.Args = args

	ret.Logs = append(ret.Logs, fmt.Sprintf("loaded [%s] from [%s] in [%s]", util.ByteSizeSI(int64(len(json))), filename, t.EndString()))
	return ret
}
