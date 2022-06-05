package project

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Save(prj *Project, logger util.Logger) error {
	fn := ConfigDir + "/project.json"
	if prj.Icon == DefaultIcon {
		prj.Icon = ""
	}
	tgtFS := s.GetFilesystem(prj)
	j := util.ToJSON(prj) + "\n"
	err := tgtFS.WriteFile(fn, []byte(j), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write project file to [%s]", fn)
	}

	if exportArgs, _ := prj.ModuleArgExport(); exportArgs != nil {
		if len(exportArgs.Config) > 0 {
			fn := ConfigDir + "/export/config.json"
			err = tgtFS.WriteFile(fn, util.ToJSONBytes(exportArgs.Config, true), filesystem.DefaultMode, true)
			if err != nil {
				return errors.Wrapf(err, "unable to write project file to [%s]", fn)
			}
		}
		for _, m := range exportArgs.Models {
			fn := ConfigDir + "/export/models/" + m.Name + ".json"
			err = tgtFS.WriteFile(fn, util.ToJSONBytes(m, true), filesystem.DefaultMode, true)
			if err != nil {
				return errors.Wrapf(err, "unable to write model file to [%s]", fn)
			}
		}
	}

	if prj.Path != "" && prj.Path != "." {
		s.addAdditional(prj.Path, logger)
	}
	return nil
}
