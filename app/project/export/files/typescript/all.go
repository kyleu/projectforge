package typescript

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
)

func All(models model.Models, enums enum.Enums, extraTypes model.Models, linebreak string) (file.Files, error) {
	var ret file.Files

	for _, e := range enums {
		x, err := EnumContent(e, nil, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing enum [%s]", e.Name)
		}
		ret = append(ret, x)
	}

	for _, m := range models {
		imps, err := tsModelImports(enums, models, extraTypes, m)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing imports for model [%s]", m.Name)
		}
		x, err := ModelContent(m, enums, imps, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing model [%s]", m.Name)
		}
		ret = append(ret, x)
	}

	return ret, nil
}
