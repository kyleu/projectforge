package load

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

func loadConfig(fs filesystem.FileLoader, exportPath string, args *metamodel.Args) error {
	exportCfgPath := util.StringFilePath(exportPath, "config.json")
	if !fs.Exists(exportCfgPath) {
		return nil
	}
	cfg, err := fs.ReadFile(exportCfgPath)
	if err != nil {
		return err
	}
	cfgMap, err := util.FromJSONMap(cfg)
	if err != nil {
		return errors.Wrap(err, "invalid export config")
	}
	args.ConfigFile = cfg
	args.Config = cfgMap
	return nil
}

func loadGroups(fs filesystem.FileLoader, exportPath string, args *metamodel.Args) error {
	groupCfgPath := util.StringFilePath(exportPath, "groups.json")
	if !fs.Exists(groupCfgPath) {
		return nil
	}
	grpFile, err := fs.ReadFile(groupCfgPath)
	if err != nil {
		return err
	}
	grps, err := util.FromJSONObj[model.Groups](grpFile)
	if err != nil {
		return errors.Wrap(err, "invalid group config")
	}
	args.GroupsFile = grpFile
	args.Groups = grps
	return nil
}

func loadTypes(fs filesystem.FileLoader, exportPath string, args *metamodel.Args) error {
	extraTypesCfgPath := util.StringFilePath(exportPath, "types.json")
	if !fs.Exists(extraTypesCfgPath) {
		return nil
	}
	extraTypesFile, err := fs.ReadFile(extraTypesCfgPath)
	if err != nil {
		return err
	}
	extraTypes, err := util.FromJSONObj[model.Models](extraTypesFile)
	if err != nil {
		return errors.Wrap(err, "invalid extraTypes config")
	}
	args.ExtraTypesFile = extraTypesFile
	args.ExtraTypes = extraTypes
	return nil
}
