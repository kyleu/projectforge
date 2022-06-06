package project

import (
	"fmt"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/export/model"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Save(prj *Project, logger util.Logger) error {
	fn := ConfigDir + "/project.json"
	if prj.Icon == DefaultIcon {
		prj.Icon = ""
	}
	tgtFS := s.GetFilesystem(prj)
	j := util.ToJSON(prj)
	err := tgtFS.WriteFile(fn, []byte(j), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write project file to [%s]", fn)
	}

	if exportArgs, _ := prj.ModuleArgExport(s, logger); exportArgs != nil {
		if len(exportArgs.Config) > 0 {
			fn := ConfigDir + "/export/config.json"
			err = tgtFS.WriteFile(fn, util.ToJSONBytes(exportArgs.Config, true), filesystem.DefaultMode, true)
			if err != nil {
				return errors.Wrapf(err, "unable to write project file to [%s]", fn)
			}
		}
		for _, m := range exportArgs.Models {
			err = s.SaveExportModel(tgtFS, m, logger)
			if err != nil {
				return err
			}
		}
	}

	if prj.Path != "" && prj.Path != "." {
		s.addAdditional(prj.Path, logger)
	}
	return nil
}

func (s *Service) SaveExportModel(fs filesystem.FileLoader, mdl *model.Model, logger util.Logger) error {
	fn := fmt.Sprintf("%s/export/models/%s.json", ConfigDir, mdl.Name)
	j := util.ToJSON(mdl) + "\n"
	err := fs.WriteFile(fn, []byte(j), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write export model file to [%s]", fn)
	}
	return nil
}

func (s *Service) DeleteExportModel(fs filesystem.FileLoader, mdl string, logger util.Logger) error {
	fn := fmt.Sprintf("%s/export/models/%s.json", ConfigDir, mdl)
	err := fs.Remove(fn, logger)
	if err != nil {
		return errors.Wrapf(err, "unable to delete export model file [%s]", fn)
	}
	return nil
}
