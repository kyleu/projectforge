package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
)

const rootKey = "root"

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
			r = rootKey
		}
		ret := NewProject(r, path)
		ret.Name = fmt.Sprintf("%s (missing)", r)
		return ret, nil
	}
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	ret := &Project{}
	err = util.FromJSON(b, &ret)
	if err != nil {
		l, r := util.StringSplitLast(path, '/', true)
		if r == "" {
			r = l
		}
		if r == "." || r == "" {
			r = rootKey
		}
		ret = NewProject(r, path)
		ret.Name = fmt.Sprintf("%s (json error)", r)
		ret.Info = &Info{}
		return ret, nil
		// return nil, errors.Wrapf(err, "can't load project from [%s]", cfgPath)
	}
	ret.Path = path

	fs := s.GetFilesystem(ret)
	if ret.Config, err = s.loadModuleConfig(fs); err != nil {
		return nil, err
	}
	if ret.Theme == nil {
		ret.Theme = theme.ThemeDefault
	}
	return ret, nil
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
