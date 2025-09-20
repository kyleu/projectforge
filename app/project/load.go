package project

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
)

const rootKey = "root"

func (s *Service) load(path string) ([]byte, *Project, error) {
	cfgPath := util.StringFilePath(path, ConfigDir, "project.json")

	rootfs, _ := filesystem.NewFileSystem(".", false, "")
	if curr, _ := rootfs.Stat(cfgPath); curr == nil {
		_, r := util.StringSplitPath(path)
		if r == "." {
			r, _ = os.Getwd()
			_, r = util.StringSplitPath(r)
			if strings.Contains(r, "/") {
				r = r[strings.LastIndex(r, "/")+1:]
			}
		}
		if r == "" {
			r = rootKey
		}
		ret := NewProject(r, path)
		ret.Name = fmt.Sprintf("%s (missing)", r)
		ret.Error = "missing [.projectforge/project.json]"
		return nil, ret, nil
	}
	b, err := rootfs.ReadFile(cfgPath)
	if err != nil {
		return nil, nil, err
	}

	ret, err := util.FromJSONObj[*Project](b)
	if err != nil {
		_, r := util.StringSplitPath(path)
		if r == "." || r == "" {
			r = rootKey
		}
		ret = NewProject(r, path)
		ret.Name = fmt.Sprintf("%s (json error)", r)
		ret.Info = &Info{}
		return b, ret, nil
		// return nil, errors.Wrapf(err, "can't load project from [%s]", cfgPath)
	}
	ret.Path = path

	fs, err := s.GetFilesystem(ret)
	if err != nil {
		return nil, nil, err
	}

	if ret.Config, err = s.loadModuleConfig(fs); err != nil {
		return b, nil, err
	}
	if ret.Theme == nil {
		ret.Theme = theme.Default
	}
	return b, ret, nil
}

func (s *Service) loadModuleConfig(fs filesystem.FileLoader) (util.ValueMap, error) {
	dbuiPath := util.StringFilePath(ConfigDir, "databaseui")
	if !fs.IsDir(dbuiPath) {
		return nil, nil
	}
	var ret util.ValueMap
	if exportCfgPath := util.StringFilePath(dbuiPath, "config.json"); fs.Exists(exportCfgPath) {
		if cfg, err := fs.ReadFile(exportCfgPath); err == nil {
			cfgMap, err := util.FromJSONMap(cfg)
			if err != nil {
				return nil, errors.Wrap(err, "invalid databaseui config")
			}
			ret = cfgMap
		}
	}
	return ret, nil
}
