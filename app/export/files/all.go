package files

import (
	"context"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/export/files/controller"
	"projectforge.dev/projectforge/app/export/files/gomodel"
	"projectforge.dev/projectforge/app/export/files/grpc"
	"projectforge.dev/projectforge/app/export/files/sql"
	"projectforge.dev/projectforge/app/export/files/svc"
	"projectforge.dev/projectforge/app/export/files/view"
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
		var calls file.Files
		var f *file.File

		fs, err := basics(m, args, addHeader)
		if err != nil {
			return nil, err
		}
		calls = append(calls, fs...)

		for _, grp := range m.GroupedColumns() {
			f, err = controller.Grouping(m, args, grp, addHeader)
			if err != nil {
				return nil, errors.Wrap(err, "can't render controller for group ["+grp.Title()+"]")
			}
			calls = append(calls, f)
		}

		if args.HasModule("migration") {
			f, err = sql.Migration(m, args, addHeader)
			if err != nil {
				return nil, errors.Wrap(err, "can't render SQL migration")
			}
			calls = append(calls, f)
		}

		fs, err = view.All(m, args, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "can't render list template")
		}
		calls = append(calls, fs...)

		if args.HasModule("grpc") && args.Config.GetStringOpt("grpcPackage") != "" {
			fs, err := grpc.GRPC(m, args, addHeader)
			if err != nil {
				return nil, err
			}
			calls = append(calls, fs...)
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
	}
	return ret, nil
}

func basics(m *model.Model, args *model.Args, addHeader bool) (file.Files, error) {
	var calls file.Files
	f, err := gomodel.Model(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render model")
	}
	calls = append(calls, f)

	if m.IsHistory() {
		f, err = gomodel.History(m, args, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "can't render History")
		}
		calls = append(calls, f)
	}

	f, err = gomodel.DTO(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render DTO")
	}
	calls = append(calls, f)

	fs, err := svc.ServiceAll(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render service")
	}
	calls = append(calls, fs...)

	f, err = controller.Controller(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render controller")
	}
	calls = append(calls, f)
	return calls, nil
}
