package inject

import (
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/model"
)

func All(args *model.Args, files file.Files, linebreak string) error {
	if args == nil {
		return nil
	}
	for _, f := range files {
		var err error
		if f.FullPath() == "app/services.go" {
			err = Services(f, args)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
