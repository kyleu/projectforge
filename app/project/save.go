package project

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/model"
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
			err = s.SaveExportModel(tgtFS, m)
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

func (s *Service) SaveExportGroups(fs filesystem.FileLoader, g model.Groups) error {
	fn := fmt.Sprintf("%s/export/groups.json", ConfigDir)
	j := util.ToJSON(g)
	err := fs.WriteFile(fn, []byte(j), filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write export groups file to [%s]", fn)
	}
	return nil
}

func (s *Service) SaveExportModel(fs filesystem.FileLoader, mdl *model.Model) error {
	if mdl.HasTag("json") {
		return nil
	}
	fn := fmt.Sprintf("%s/export/models/%s.json", ConfigDir, mdl.Name)
	j := util.ToJSONBytes(mdl, true)
	err := fs.WriteFile(fn, j, filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write export model file to [%s]", fn)
	}
	return nil
}

func (s *Service) DeleteExportModel(fs filesystem.FileLoader, mdl string, logger util.Logger) error {
	fn := fmt.Sprintf("%s/export/models/%s.json", ConfigDir, mdl)
	if err := fs.Remove(fn, logger); err != nil {
		return errors.Wrapf(err, "unable to delete export model file [%s]", fn)
	}
	return nil
}

func (s *Service) SaveExportEnum(fs filesystem.FileLoader, e *enum.Enum) error {
	fn := fmt.Sprintf("%s/export/enums/%s.json", ConfigDir, e.Name)
	j := util.ToJSONBytes(e, true)
	err := fs.WriteFile(fn, j, filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write export enum file to [%s]", fn)
	}
	return nil
}

func (s *Service) DeleteExportEnum(fs filesystem.FileLoader, e string, logger util.Logger) error {
	fn := fmt.Sprintf("%s/export/enums/%s.json", ConfigDir, e)
	if err := fs.Remove(fn, logger); err != nil {
		return errors.Wrapf(err, "unable to delete export enum file [%s]", fn)
	}
	return nil
}
