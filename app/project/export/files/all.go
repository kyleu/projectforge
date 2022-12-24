package files

import (
	"context"
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/filesystem"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files/controller"
	"projectforge.dev/projectforge/app/project/export/files/goenum"
	"projectforge.dev/projectforge/app/project/export/files/sql"
	"projectforge.dev/projectforge/app/project/export/model"
	"projectforge.dev/projectforge/app/util"
)

func All(ctx context.Context, p *project.Project, args *model.Args, addHeader bool, logger util.Logger) (file.Files, error) {
	if err := args.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid export arguments")
	}
	ret := make(file.Files, 0, (len(args.Models)*10)+len(args.Enums))

	for _, e := range args.Enums {
		call, err := goenum.Enum(e, args, addHeader)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing enum [%s]", e.Name)
		}
		ret = append(ret, call)
	}

	for _, m := range args.Models {
		calls, err := ModelAll(m, p, args, addHeader)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing model [%s]", m.Name)
		}
		ret = append(ret, calls...)
	}

	x, err := controller.Routes(args, addHeader)
	if err != nil {
		return nil, err
	}
	ret = append(ret, x)

	x, err = controller.Menu(args, addHeader)
	if err != nil {
		return nil, err
	}
	ret = append(ret, x)

	if args.HasModule("search") {
		x, err = controller.Search(args, addHeader)
		if err != nil {
			return nil, err
		}
		ret = append(ret, x)
	}

	if args.HasModule("datadog") {
		svc := ServiceDefinition(p)
		f := file.NewFile("doc/service.json", filesystem.DefaultMode, util.ToJSONBytes(svc, true), false, logger)
		ret = append(ret, f)
	}
	if args.HasModule("migration") {
		f, err := sql.MigrationAll(args.Models.Sorted(), args.Enums, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL \"all\" migration")
		}
		ret = append(ret, f)
	}
	if args.Models.HasSeedData() {
		f, err := sql.SeedDataAll(args.Models)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL \"all\" migration")
		}
		ret = append(ret, f)
	}
	if len(args.Enums) > 0 {
		f, err := sql.Types(args.Enums, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL types")
		}
		ret = append(ret, f)
	}
	return ret, nil
}
