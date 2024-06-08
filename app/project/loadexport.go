package project

import (
	"cmp"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/data"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) loadExportArgs(fs filesystem.FileLoader, acronyms []string, logger util.Logger) (*model.Args, error) {
	args := &model.Args{Config: util.ValueMap{}}
	exportPath := filepath.Join(ConfigDir, "export")
	if !fs.IsDir(exportPath) {
		return args, nil
	}
	if exportCfgPath := filepath.Join(exportPath, "config.json"); fs.Exists(exportCfgPath) {
		if cfg, err := fs.ReadFile(exportCfgPath); err == nil {
			cfgMap := util.ValueMap{}
			err = util.FromJSON(cfg, &cfgMap)
			if err != nil {
				return nil, errors.Wrap(err, "invalid export config")
			}
			args.ConfigFile = cfg
			args.Config = cfgMap
		}
	}
	if groupCfgPath := filepath.Join(exportPath, "groups.json"); fs.Exists(groupCfgPath) {
		if grpFile, err := fs.ReadFile(groupCfgPath); err == nil {
			grps := model.Groups{}
			err = util.FromJSON(grpFile, &grps)
			if err != nil {
				return nil, errors.Wrap(err, "invalid group config")
			}
			args.GroupsFile = grpFile
			args.Groups = grps
		}
	}

	enumFiles, enums, err := getEnums(exportPath, fs, logger)
	if err != nil {
		return nil, err
	}
	args.EnumFiles = enumFiles
	args.Enums = enums

	explicitModelFiles, explicitModels, err := getModels(exportPath, fs, logger)
	if err != nil {
		return nil, err
	}
	args.ModelFiles = explicitModelFiles
	args.Models = append(args.Models, explicitModels...)

	jsonModelFiles, jsonModels, err := getJSONModels(args.Config, args.Groups, fs, logger)
	if err != nil {
		return nil, err
	}
	for k, v := range jsonModelFiles {
		args.ModelFiles[k] = v
	}
	args.Models = append(args.Models, jsonModels...)
	if len(args.Models) > 0 {
		slices.SortFunc(args.Models, func(l *model.Model, r *model.Model) int {
			return cmp.Compare(l.Name, r.Name)
		})
	}
	args.ApplyAcronyms(acronyms...)
	return args, nil
}

func getEnums(exportPath string, fs filesystem.FileLoader, logger util.Logger) (map[string]json.RawMessage, enum.Enums, error) {
	enumsPath := filepath.Join(exportPath, "enums")
	if !fs.IsDir(enumsPath) {
		return nil, nil, nil
	}
	enumNames := fs.ListJSON(enumsPath, nil, false, logger)
	enums := make(enum.Enums, 0, len(enumNames))
	enumFiles := make(map[string]json.RawMessage, len(enumNames))
	for _, enumName := range enumNames {
		fn := filepath.Join(enumsPath, enumName)
		content, err := fs.ReadFile(fn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read export enum file from [%s]", fn)
		}
		m := &enum.Enum{}
		err = util.FromJSON(content, m)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read export enum JSON from [%s]", fn)
		}
		enums = append(enums, m)
		enumFiles[m.Name] = content
	}
	return enumFiles, enums, nil
}

func getModels(exportPath string, fs filesystem.FileLoader, logger util.Logger) (map[string]json.RawMessage, model.Models, error) {
	modelsPath := filepath.Join(exportPath, "models")
	if !fs.IsDir(modelsPath) {
		return nil, nil, nil
	}
	modelNames := fs.ListJSON(modelsPath, nil, false, logger)
	models := make(model.Models, 0, len(modelNames))
	modelFiles := make(map[string]json.RawMessage, len(modelNames))
	for _, modelName := range modelNames {
		fn := filepath.Join(modelsPath, modelName)
		content, err := fs.ReadFile(fn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read export model file from [%s]", fn)
		}
		m := &model.Model{}
		err = util.FromJSON(content, m)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read export model JSON from [%s]", fn)
		}
		modelFiles[m.Name] = content
		models = append(models, m)
	}
	return modelFiles, models, nil
}

func getJSONModels(cfg util.ValueMap, groups model.Groups, fs filesystem.FileLoader, logger util.Logger) (map[string]json.RawMessage, model.Models, error) {
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
		fn := path.Join(pth, jsonFile)
		x, err := fs.ReadFile(fn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "can't read file [%s]", fn)
		}
		df := &data.File{}
		err = util.FromJSON(x, df)
		if err != nil {
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
