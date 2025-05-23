package load

import (
	"cmp"
	"slices"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/util"
)

func ExportArgs(fs filesystem.FileLoader, cfgDir string, acronyms []string, logger util.Logger) (*metamodel.Args, error) {
	args := &metamodel.Args{Config: util.ValueMap{}}
	exportPath := util.StringFilePath(cfgDir, "export")
	if !fs.IsDir(exportPath) {
		return args, nil
	}

	if err := loadConfig(fs, exportPath, args); err != nil {
		return nil, err
	}
	if err := loadGroups(fs, exportPath, args); err != nil {
		return nil, err
	}
	if err := loadTypes(fs, exportPath, args); err != nil {
		return nil, err
	}

	enumFiles, enums, err := Enums(exportPath, fs, logger)
	if err != nil {
		return nil, err
	}
	args.EnumFiles = enumFiles
	args.Enums = enums

	explicitModelFiles, explicitModels, err := LoadModels(exportPath, fs, logger)
	if err != nil {
		return nil, err
	}
	args.ModelFiles = explicitModelFiles
	args.Models = append(args.Models, explicitModels...)

	jsonModelFiles, jsonModels, err := LoadJSONModels(args.Config, args.Groups, fs, logger)
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
