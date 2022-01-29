package files

import (
	"github.com/pkg/errors"

	"github.com/kyleu/projectforge/app/export/files/controller"
	"github.com/kyleu/projectforge/app/export/files/gomodel"
	"github.com/kyleu/projectforge/app/export/files/grpc"
	"github.com/kyleu/projectforge/app/export/files/sql"
	"github.com/kyleu/projectforge/app/export/files/svc"
	"github.com/kyleu/projectforge/app/export/files/view"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
)

func All(args *model.Args, addHeader bool) (file.Files, error) {
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

		if args.HasModule("grpc") {
			fs, err := grpc.GRPC(m, args, addHeader)
			if err != nil {
				return nil, err
			}
			calls = append(calls, fs...)
		}

		ret = append(ret, calls...)
	}

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
