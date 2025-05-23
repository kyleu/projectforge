package load

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/data"
	"projectforge.dev/projectforge/app/util"
)

func LoadModels(exportPath string, fs filesystem.FileLoader, logger util.Logger) (map[string]json.RawMessage, model.Models, error) {
	modelsPath := util.StringFilePath(exportPath, "models")
	if !fs.IsDir(modelsPath) {
		return nil, nil, nil
	}
	modelNames := fs.ListJSON(modelsPath, nil, false, logger)
	models := make(model.Models, 0, len(modelNames))
	modelFiles := make(map[string]json.RawMessage, len(modelNames))
	for _, modelName := range modelNames {
		fn := util.StringFilePath(modelsPath, modelName)
		content, err := fs.ReadFile(fn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read export model file from [%s]", fn)
		}
		m, err := util.FromJSONObj[*model.Model](content)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read export model JSON from [%s]", fn)
		}
		modelFiles[m.Name] = content
		models = append(models, m)
	}
	return modelFiles, models, nil
}

func LoadJSONModels(cfg util.ValueMap, groups model.Groups, fs filesystem.FileLoader, logger util.Logger) (map[string]json.RawMessage, model.Models, error) {
	if cfg == nil {
		return nil, nil, nil
	}
	jsonCfg, ok := cfg["jsonExport"]
	if !ok {
		return nil, nil, nil
	}
	pth := fmt.Sprint(jsonCfg)
	if pth == "" {
		return nil, nil, nil
	}
	jsonFiles := fs.ListJSON(pth, nil, false, logger)
	jsonModelFiles := make(map[string]json.RawMessage, len(jsonFiles))
	jsonModels := make(model.Models, 0, len(jsonFiles))
	for idx, jsonFile := range jsonFiles {
		if !strings.HasSuffix(jsonFile, util.ExtJSON) {
			continue
		}
		if strings.Contains(jsonFile, ".min.") {
			continue
		}
		fn := util.StringFilePath(pth, jsonFile)
		x, err := fs.ReadFile(fn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "can't read file [%s]", fn)
		}
		df := &data.File{}
		if err = util.FromJSONStrict(x, df); err != nil {
			return nil, nil, errors.Wrapf(err, "unable to parse JSON data file at [%s]", fn)
		}
		mdl, err := df.ToModel(idx, groups)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to convert file at [%s]", fn)
		}
		jsonModelFiles[mdl.Name] = x
		jsonModels = append(jsonModels, mdl)
	}
	return jsonModelFiles, jsonModels, nil
}
