package files

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/export/files/controller"
	"projectforge.dev/projectforge/app/export/files/sql"
	"projectforge.dev/projectforge/app/export/model"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/util"
)

func All(ctx context.Context, args *model.Args, addHeader bool, logger util.Logger) (file.Files, error) {
	if err := args.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid export arguments")
	}
	ret := make(file.Files, 0, len(args.Models)*10)
	for _, m := range args.Models {
		calls, err := modelAll(m, args, addHeader)
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

	x, err = controller.Search(args, addHeader)
	if err != nil {
		return nil, err
	}
	ret = append(ret, x)

	if args.HasModule("migration") {
		f, err := sql.MigrationAll(args.Models, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL \"all\" migration")
		}
		ret = append(ret, f)
		if args.Models.HasSeedData() {
			f, err = sql.SeedDataAll(args.Models)
			if err != nil {
				return nil, errors.Wrap(err, "can't render SQL \"all\" migration")
			}
			ret = append(ret, f)
		}
	}
	return ret, nil
}
