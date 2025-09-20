package files

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files/goenum"
	"projectforge.dev/projectforge/app/project/export/files/sql"
)

func EnumAll(p *project.Project, linebreak string, enums enum.Enums, database string) (file.Files, error) {
	var ret file.Files
	for _, e := range enums {
		call, err := goenum.Enum(e, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing enum [%s]", e.Name)
		}
		ret = append(ret, call)
	}

	if len(enums.WithDatabase()) > 0 && p.HasModule("migration") {
		f, err := sql.Types(enums.WithDatabase(), linebreak, database)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL types")
		}
		ret = append(ret, f)
	}
	return ret, nil
}
