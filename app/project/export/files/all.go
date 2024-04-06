package files

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files/controller"
	"projectforge.dev/projectforge/app/project/export/files/goenum"
	"projectforge.dev/projectforge/app/project/export/files/gql"
	"projectforge.dev/projectforge/app/project/export/files/script"
	"projectforge.dev/projectforge/app/project/export/files/sql"
	"projectforge.dev/projectforge/app/project/export/files/svc"
	"projectforge.dev/projectforge/app/project/export/model"
)

func All(p *project.Project, args *model.Args, addHeader bool, linebreak string) (file.Files, error) {
	if err := args.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid export arguments")
	}
	ret := make(file.Files, 0, (len(args.Models)*10)+len(args.Enums))

	for _, e := range args.Enums {
		call, err := goenum.Enum(e, addHeader, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing enum [%s]", e.Name)
		}
		ret = append(ret, call)
	}

	for _, m := range args.Models {
		calls, err := ModelAll(m, p, args, addHeader, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing model [%s]", m.Name)
		}
		ret = append(ret, calls...)
	}

	x, err := svc.Services(args, addHeader, linebreak)
	if err != nil {
		return nil, err
	}
	ret = append(ret, x)

	x, err = controller.Routes(args, addHeader, linebreak)
	if err != nil {
		return nil, err
	}
	ret = append(ret, x)

	x, err = controller.Menu(args, addHeader, linebreak)
	if err != nil {
		return nil, err
	}
	ret = append(ret, x)

	if args.HasModule("search") {
		x, err = controller.Search(args, addHeader, linebreak)
		if err != nil {
			return nil, err
		}
		ret = append(ret, x)
	}

	if args.HasModule("graphql") {
		x, err = gql.All(args.Models, args.Enums, addHeader, linebreak)
		if err != nil {
			return nil, err
		}
		ret = append(ret, x)
	}

	if args.HasModule("migration") {
		migModels := args.Models.WithoutTag("external").Sorted()
		f, err := sql.MigrationAll(migModels, args.Enums, addHeader, linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL \"all\" migration")
		}
		ret = append(ret, f)
	}

	if args.HasModule("notebook") {
		x, err := script.NotebookScript(p, args, addHeader, linebreak)
		if err != nil {
			return nil, err
		}
		ret = append(ret, x)
	}

	if args.Models.HasSeedData() {
		f, err := sql.SeedDataAll(args.Models, linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL \"all\" migration")
		}
		ret = append(ret, f)
	}
	if len(args.Enums) > 0 {
		f, err := sql.Types(args.Enums, addHeader, linebreak, args.Database)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL types")
		}
		ret = append(ret, f)
	}
	return ret, nil
}
