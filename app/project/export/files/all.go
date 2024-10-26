package files

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/enum"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files/controller"
	"projectforge.dev/projectforge/app/project/export/files/goenum"
	"projectforge.dev/projectforge/app/project/export/files/gql"
	"projectforge.dev/projectforge/app/project/export/files/script"
	"projectforge.dev/projectforge/app/project/export/files/sql"
	"projectforge.dev/projectforge/app/project/export/files/svc"
)

func All(p *project.Project, linebreak string) (file.Files, error) {
	if p.ExportArgs == nil {
		return nil, errors.New("export arguments aren't loaded")
	}
	args := p.ExportArgs
	if err := args.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid export arguments")
	}
	ret := make(file.Files, 0, (len(args.Models)*10)+len(args.Enums))

	enumRet, err := allEnum(p, linebreak, args.Enums, args.Database)
	if err != nil {
		return nil, err
	}
	ret = append(ret, enumRet...)

	for _, m := range args.Models {
		calls, e := ModelAll(m, p, args, linebreak)
		if e != nil {
			return nil, errors.Wrapf(e, "error processing model [%s]", m.Name)
		}
		ret = append(ret, calls...)
	}

	x, err := svc.Services(args, linebreak)
	if err != nil {
		return nil, err
	}
	ret = append(ret, x)

	x, err = controller.Routes(args, linebreak)
	if err != nil {
		return nil, err
	}
	ret = append(ret, x)

	x, err = controller.Menu(args, linebreak)
	if err != nil {
		return nil, err
	}
	ret = append(ret, x)

	extraRet, err := extraFiles(p, linebreak, args)
	if err != nil {
		return nil, err
	}
	ret = append(ret, extraRet...)

	return ret, nil
}

func allEnum(p *project.Project, linebreak string, enums enum.Enums, database string) (file.Files, error) {
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

func extraFiles(p *project.Project, linebreak string, args *model.Args) (file.Files, error) {
	var ret file.Files

	if args.HasModule("search") {
		x, err := controller.Search(args, linebreak)
		if err != nil {
			return nil, err
		}
		ret = append(ret, x)
	}

	if args.HasModule("graphql") {
		m, e := args.Models.WithDatabase(), args.Enums.WithDatabase()
		if len(m) > 0 || len(e) > 0 {
			x, err := gql.All(m, e, linebreak)
			if err != nil {
				return nil, err
			}
			ret = append(ret, x)
		}
	}

	if args.HasModule("notebook") {
		x, err := script.NotebookScript(p, args.Models, linebreak)
		if err != nil {
			return nil, err
		}
		ret = append(ret, x)
	}

	if args.HasModule("migration") {
		migModels := args.Models.WithoutTag("external").WithDatabase().Sorted()
		migEnums := args.Enums.WithDatabase()
		f, err := sql.MigrationAll(migModels, migEnums, linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL \"all\" migration")
		}
		ret = append(ret, f)
	}

	if len(args.Models.WithDatabase()) > 0 || len(args.Enums.WithDatabase()) > 0 {
		f, err := sql.SeedDataAll(args.Models.WithDatabase(), linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL \"all\" migration")
		}
		ret = append(ret, f)
	}

	return ret, nil
}
