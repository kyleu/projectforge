package project

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) load(path string, logger util.Logger) (*Project, error) {
	cfgPath := filepath.Join(path, ConfigDir, "project.json")
	if curr, _ := os.Stat(cfgPath); curr == nil {
		l, r := util.StringSplitLast(path, '/', true)
		if r == "" {
			r = l
		}
		if r == "." {
			r, _ = os.Getwd()
			if strings.Contains(r, "/") {
				r = r[strings.LastIndex(r, "/")+1:]
			}
		}
		if r == "" {
			r = "root"
		}
		ret := NewProject(r, path)
		ret.Name = r + " (missing)"
		return ret, nil
	}
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	ret := &Project{}
	err = util.FromJSON(b, &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "can't load project from [%s]", cfgPath)
	}
	ret.Path = path

	fs := s.GetFilesystem(ret)
	if ret.ExportArgs, err = s.loadExportArgs(fs, logger); err != nil {
		return nil, err
	}

	if ret.Config, err = s.loadModuleConfig(fs); err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *Service) loadExportArgs(fs filesystem.FileLoader, logger util.Logger) (*model.Args, error) {
	exportPath := filepath.Join(ConfigDir, "export")
	if !fs.IsDir(exportPath) {
		return &model.Args{}, nil
	}
	args := &model.Args{}
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
	if modelsPath := filepath.Join(exportPath, "models"); fs.IsDir(modelsPath) {
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

		args.Models = models.Sorted()
	}
	return args, nil
}

func (s *Service) loadModuleConfig(fs filesystem.FileLoader) (util.ValueMap, error) {
	dbuiPath := filepath.Join(ConfigDir, "databaseui")
	if !fs.IsDir(dbuiPath) {
		return nil, nil
	}
	var ret util.ValueMap
	if exportCfgPath := filepath.Join(dbuiPath, "config.json"); fs.Exists(exportCfgPath) {
		if cfg, err := fs.ReadFile(exportCfgPath); err == nil {
			cfgMap := util.ValueMap{}
			err = util.FromJSON(cfg, &cfgMap)
			if err != nil {
				return nil, errors.Wrap(err, "invalid databaseui config")
			}
			ret = cfgMap
		}
	}
	return ret, nil
}
