package files

import (
	"github.com/kyleu/projectforge/app/export/files/controller"
	"github.com/kyleu/projectforge/app/export/files/gomodel"
	"github.com/kyleu/projectforge/app/export/files/grpc"
	"github.com/kyleu/projectforge/app/export/files/sql"
	"github.com/kyleu/projectforge/app/export/files/svc"
	"github.com/kyleu/projectforge/app/export/files/view"
	"github.com/kyleu/projectforge/app/export/model"
	"github.com/kyleu/projectforge/app/file"
	"github.com/pkg/errors"
)

func All(args *model.Args) (file.Files, error) {
	ret := make(file.Files, 0, len(args.Models)*10)
	for _, m := range args.Models {
		var calls file.Files
		var f *file.File
		err := m.Validate()
		if err != nil {
			return nil, errors.Wrap(err, "invalid model ["+m.Name+"]")
		}

		f, err = gomodel.Model(m, args)
		if err != nil {
			return nil, errors.Wrap(err, "can't render model")
		}
		calls = append(calls, f)

		f, err = gomodel.DTO(m, args)
		if err != nil {
			return nil, errors.Wrap(err, "can't render DTO")
		}
		calls = append(calls, f)

		fs, err := svc.ServiceAll(m, args)
		if err != nil {
			return nil, errors.Wrap(err, "can't render service")
		}
		calls = append(calls, fs...)

		f, err = controller.Controller(m, args)
		if err != nil {
			return nil, errors.Wrap(err, "can't render controller")
		}
		calls = append(calls, f)

		for _, grp := range m.GroupedColumns() {
			f, err = controller.Grouping(m, args, grp)
			if err != nil {
				return nil, errors.Wrap(err, "can't render controller for group ["+grp.Title()+"]")
			}
			calls = append(calls, f)
		}

		if args.HasModule("migration") {
			f, err = sql.Migration(m, args)
			if err != nil {
				return nil, errors.Wrap(err, "can't render SQL migration")
			}
			calls = append(calls, f)
		}

		fs, err = view.All(m, args)
		if err != nil {
			return nil, errors.Wrap(err, "can't render list template")
		}
		calls = append(calls, fs...)

		if args.HasModule("grpc") {
			fs, err := grpc.GRPC(m, args)
			if err != nil {
				return nil, err
			}
			calls = append(calls, fs...)
		}

		ret = append(ret, calls...)
	}

	if args.HasModule("migration") {
		f, err := sql.MigrationAll(args.Models)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL \"all\" migration")
		}
		ret = append(ret, f)
	}
	return ret, nil
}
