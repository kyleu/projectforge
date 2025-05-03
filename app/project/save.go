package project

import (
	"fmt"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/theme"
	"projectforge.dev/projectforge/app/util"
)

func (s *Service) Save(prj *Project, logger util.Logger) error {
	fn := ConfigDir + "/project.json"
	if prj.Icon == DefaultIcon {
		prj.Icon = ""
	}
	if prj.Exec == prj.Key {
		prj.Exec = ""
	}
	if prj.Theme != nil && prj.Theme.Equals(theme.Default) {
		prj.Theme = nil
		defer func() {
			prj.Theme = theme.Default
		}()
	}
	tgtFS, err := s.GetFilesystem(prj)
	if err != nil {
		return err
	}
	err = tgtFS.WriteJSONFile(fn, prj, filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write project file to [%s]", fn)
	}

	err = prj.ModuleArgExport(s, logger)
	if err != nil {
		return errors.Wrapf(err, "unable to load [%s] export configuration", prj.Key)
	}
	if prj.ExportArgs != nil {
		if err = s.SaveExportArgs(tgtFS, prj.ExportArgs); err != nil {
			return errors.Wrapf(err, "unable to save export args for [%s]", prj.Key)
		}
	}

	if prj.Path != "" && prj.Path != "." {
		fs, _ := filesystem.NewFileSystem(".", false, "")
		s.addAdditional(prj.Path, fs, logger)
	}
	return nil
}

func (s *Service) SaveExportArgs(fs filesystem.FileLoader, args *metamodel.Args) error {
	if len(args.Config) > 0 {
		if err := s.SaveExportConfig(fs, args.Config); err != nil {
			return err
		}
	}
	if len(args.Groups) > 0 {
		if err := s.SaveExportGroups(fs, args.Groups); err != nil {
			return err
		}
	}
	for _, m := range args.Models {
		if err := s.SaveExportModel(fs, m); err != nil {
			return err
		}
	}
	for _, e := range args.Enums {
		if err := s.SaveExportEnum(fs, e); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) SaveExportConfig(fs filesystem.FileLoader, cfg util.ValueMap) error {
	fn := ConfigDir + "/export/config.json"
	err := fs.WriteJSONFile(fn, cfg, filesystem.DefaultMode, true)
	if err != nil {
		return errors.Wrapf(err, "unable to write export config file to [%s]", fn)
	}
	return nil
}

func (s *Service) SaveExportGroups(fs filesystem.FileLoader, g model.Groups) error {
	fn := fmt.Sprintf("%s/export/groups.json", ConfigDir)
	err := fs.WriteJSONFile(fn, g, filesystem.DefaultMode, true)
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
	err := fs.WriteJSONFile(fn, mdl, filesystem.DefaultMode, true)
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
	err := fs.WriteJSONFile(fn, e, filesystem.DefaultMode, true)
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
