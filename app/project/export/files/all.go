package files

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files/controller"
	"projectforge.dev/projectforge/app/project/export/files/gql"
	"projectforge.dev/projectforge/app/project/export/files/script"
	"projectforge.dev/projectforge/app/project/export/files/sql"
	"projectforge.dev/projectforge/app/project/export/files/svc"
	"projectforge.dev/projectforge/app/project/export/files/typescript"
)

func All(p *project.Project, linebreak string) (file.Files, error) {
	if p.ExportArgs == nil {
		return nil, errors.New("export arguments aren't loaded")
	}
	args := p.ExportArgs
	if err := args.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid export arguments")
	}
	ret := make(file.Files, 0, (len(args.Enums) + len(args.Events) + len(args.Models)*10))

	enumFiles, err := EnumAll(p, linebreak, args.Enums, args.Database)
	if err != nil {
		return nil, err
	}
	ret = append(ret, enumFiles...)

	for _, m := range args.Models {
		files, e := ModelAll(m, p, linebreak)
		if e != nil {
			return nil, errors.Wrapf(e, "error processing model [%s]", m.Name)
		}
		ret = append(ret, files...)
	}

	files, err := EventAll(p, args.Events, linebreak)
	if err != nil {
		return nil, err
	}
	ret = append(ret, files...)

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

	extraFiles, err := extraFiles(p, linebreak, args)
	if err != nil {
		return nil, err
	}
	ret = append(ret, extraFiles...)

	return ret, nil
}

func extraFiles(p *project.Project, linebreak string, args *metamodel.Args) (file.Files, error) {
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
		if wd := args.Models.WithDatabase(); len(wd) > 0 {
			x, err := script.NotebookScript(p, wd, linebreak)
			if err != nil {
				return nil, err
			}
			ret = append(ret, x)
		}
	}

	if args.HasModule("migration") {
		migModels := args.Models.WithoutTag("external").WithDatabase().Sorted()
		migEnums := args.Enums.WithDatabase()
		f, err := sql.MigrationAll(migModels, migEnums, args.HasModule("audit"), linebreak)
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

	if en, evt, mdl := args.Enums.WithTypeScript(), args.Events, args.Models.WithTypeScript(); len(en) > 0 || len(evt) > 0 || len(mdl) > 0 {
		files, err := typescript.All(args, linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render TypeScript output")
		}
		ret = append(ret, files...)
	}

	return ret, nil
}
