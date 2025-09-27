package typescript

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
)

func All(args *metamodel.Args, linebreak string) (file.Files, error) {
	var ret file.Files

	for _, e := range args.Enums.WithTypeScript() {
		x, err := EnumContent(e, nil, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing enum [%s]", e.Name)
		}
		ret = append(ret, x)
	}

	for _, evt := range args.Events.WithTypeScript() {
		imps, err := tsModelImports(args, evt.Columns.NotDerived(), evt)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing imports for model [%s]", evt.Name)
		}
		x, err := EventContent(evt, args.Enums, imps, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing model [%s]", evt.Name)
		}
		ret = append(ret, x)
	}

	for _, m := range args.Models.WithTypeScript() {
		imps, err := tsModelImports(args, m.Columns, m)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing imports for model [%s]", m.Name)
		}
		x, err := ModelContent(m, args.Enums, imps, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing model [%s]", m.Name)
		}
		ret = append(ret, x)
	}

	return ret, nil
}
