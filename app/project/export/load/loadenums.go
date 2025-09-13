package load

import (
	"encoding/json"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/util"
)

func LoadEnums(exportPath string, fs filesystem.FileLoader, logger util.Logger) (map[string]json.RawMessage, enum.Enums, error) {
	enumsPath := util.StringFilePath(exportPath, "enums")
	if !fs.IsDir(enumsPath) {
		return nil, nil, nil
	}
	enumNames := fs.ListJSON(enumsPath, nil, false, logger)
	enums := make(enum.Enums, 0, len(enumNames))
	enumFiles := make(map[string]json.RawMessage, len(enumNames))
	for _, enumName := range enumNames {
		fn := util.StringFilePath(enumsPath, enumName)
		content, err := fs.ReadFile(fn)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read export enum file from [%s]", fn)
		}
		e, err := util.FromJSONObj[*enum.Enum](content)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to read export enum JSON from [%s]", fn)
		}
		enums = append(enums, e)
		enumFiles[e.Name] = content
	}
	return enumFiles, enums, nil
}
