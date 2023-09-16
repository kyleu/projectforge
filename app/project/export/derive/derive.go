package derive

import (
	"os"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func Derive(name string, pkg string, content string, logger util.Logger) Result {
	ret := Result{}
	var js any
	if err := util.FromJSON([]byte(content), &js); err == nil {
		return deriveJSON(name, pkg, js, ret, logger)
	}
	if st, _ := os.Stat(content); st != nil {
		pth := content
		if st.IsDir() {
			dirFS, _ := filesystem.NewFileSystem(pth, true, "")
			return deriveDirectory(pth, dirFS, ret, logger)
		}
		b, err := os.ReadFile(pth)
		if err != nil {
			return ret.AddErrorSection(pth, err)
		}
		return deriveFileContent(pth, string(b), ret, logger)
	}
	return ret.AddSection("error", (&Section{}).AddError(errors.Errorf("content isn't valid JSON or a valid path [%s]", content)))
}

func deriveFileContent(fn string, content string, ret Result, logger util.Logger) Result {
	s := &Section{}
	s.AddLog("loading file [%s] (%s)", fn, util.ByteSizeSI(int64(len(content))))

	return ret.AddSection(fn, s)
}

func deriveDirectory(fn string, fs filesystem.FileLoader, ret Result, logger util.Logger) Result {
	s := &Section{}
	files, err := fs.ListFilesRecursive(".", nil, logger)
	if err != nil {
		return ret.AddErrorSection(fn, err)
	}
	s.AddLog("loading directory [%s] (%s)", fn, util.StringPlural(len(files), "file"))
	return ret.AddSection(fn, s)
}

func deriveJSON(name string, pkg string, js any, ret Result, logger util.Logger) Result {
	switch t := js.(type) {
	case map[string]any:
		return deriveJSONObject(name, pkg, t, ret)
	default:
		return ret.AddErrorSection("json", errors.Errorf("unhandled type [%T]: %v", js, t))
	}
}

func deriveJSONObject(name string, pkg string, m map[string]any, ret Result) Result {
	if pkg == "" {
		pkg = name
	}
	cols, seedData, err := deriveExtractColumns("", m)
	if err != nil {
		return ret.AddErrorSection("json", err)
	}
	return ret.AddSection(name, (&Section{}).AddModel(&model.Model{Name: name, Package: pkg, Icon: "star", Columns: cols, SeedData: seedData}))
}

func deriveExtractColumns(prefix string, childMap util.ValueMap) (model.Columns, [][]any, error) {
	var ret model.Columns
	var sd [][]any
	for k, v := range childMap {
		col := &model.Column{Name: k, PK: k == "id"}
		switch t := v.(type) {
		case float64:
			if t == float64(int64(t)) {
				col.Type = types.NewInt(64)
			} else {
				col.Type = types.NewFloat(64)
			}
		case string:
			col.Type = types.NewString()
		case map[string]any:
			newCols, newSeed, err := deriveExtractColumns(k, t)
			if err != nil {
				return nil, nil, err
			}
			ret = append(ret, newCols...)
			sd = append(sd, newSeed...)
		default:
			return nil, nil, errors.Errorf("unhandled type [%T]: %v", v, t)
		}
		if col.Type != nil {
			ret = append(ret, col)
		}
	}
	return ret, sd, nil
}
