package project

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"projectforge.dev/projectforge/app/project/export/data"
	"projectforge.dev/projectforge/app/project/export/model"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) loadExportArgs(fs filesystem.FileLoader, logger util.Logger) (*model.Args, error) {
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
			args.Config = cfgMap
		}
	}
	if groupCfgPath := filepath.Join(exportPath, "groups.json"); fs.Exists(groupCfgPath) {
		if cfg, err := fs.ReadFile(groupCfgPath); err == nil {
			grps := model.Groups{}
			err = util.FromJSON(cfg, &grps)
			if err != nil {
				return nil, errors.Wrap(err, "invalid group config")
			}
			args.Groups = grps
		}
	}
	explicitModels, err := getModels(exportPath, fs, logger)
	if err != nil {
		return nil, err
	}
	args.Models = append(args.Models, explicitModels...)

	jsonModels, err := getJSONModels(args.Config, args.Groups, fs, logger)
	if err != nil {
		return nil, err
	}
	args.Models = append(args.Models, jsonModels...)

	if len(args.Models) > 0 {
		slices.SortFunc(args.Models, func(l *model.Model, r *model.Model) bool {
			return l.Name < r.Name
		})
	}

	return args, nil
}

func getModels(exportPath string, fs filesystem.FileLoader, logger util.Logger) (model.Models, error) {
	modelsPath := filepath.Join(exportPath, "models")
	if !fs.IsDir(modelsPath) {
		return nil, nil
	}
	modelNames := fs.ListJSON(modelsPath, nil, false, logger)
	models := model.Models{}
	for _, modelName := range modelNames {
		fn := filepath.Join(modelsPath, modelName)
		content, err := fs.ReadFile(fn)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read export model file from [%s]", fn)
		}
		m := &model.Model{}
		err = util.FromJSON(content, m)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read export model JSON from [%s]", fn)
		}
		models = append(models, m)
	}
	return models, nil
}

func getJSONModels(cfg util.ValueMap, groups model.Groups, fs filesystem.FileLoader, logger util.Logger) (model.Models, error) {
	if cfg == nil {
		return nil, nil
	}
	jsonCfg, ok := cfg["jsonExport"]
	if !ok {
		return nil, nil
	}
	pth := fmt.Sprint(jsonCfg)
	if pth == "" {
		return nil, nil
	}
	jsonFiles := fs.ListJSON(pth, nil, false, logger)
	jsonModels := make(model.Models, 0, len(jsonFiles))
	for idx, jsonFile := range jsonFiles {
		if !strings.HasSuffix(jsonFile, ".json") {
			continue
		}
		if strings.Contains(jsonFile, ".min.") {
			continue
		}
		fn := path.Join(pth, jsonFile)
		x, err := fs.ReadFile(fn)
		if err != nil {
			return nil, errors.Wrapf(err, "can't read file [%s]", fn)
		}
		df := &data.File{}
		err = util.FromJSON(x, df)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to parse JSON data file at [%s]", fn)
		}
		mdl, err := df.ToModel(idx, groups)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to convert file at [%s]", fn)
		}
		jsonModels = append(jsonModels, mdl)
	}
	return jsonModels, nil
}
