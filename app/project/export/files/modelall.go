package files

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/file"
	controller2 "projectforge.dev/projectforge/app/project/export/files/controller"
	gomodel2 "projectforge.dev/projectforge/app/project/export/files/gomodel"
	"projectforge.dev/projectforge/app/project/export/files/grpc"
	sql2 "projectforge.dev/projectforge/app/project/export/files/sql"
	"projectforge.dev/projectforge/app/project/export/files/svc"
	"projectforge.dev/projectforge/app/project/export/files/view"
	model2 "projectforge.dev/projectforge/app/project/export/model"
)

func ModelAll(m *model2.Model, args *model2.Args, addHeader bool) (file.Files, error) {
	var calls file.Files
	var f *file.File

	fs, err := basics(m, args, addHeader)
	if err != nil {
		return nil, err
	}
	calls = append(calls, fs...)

	for _, grp := range m.GroupedColumns() {
		f, err = controller2.Grouping(m, args, grp, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "can't render controller for group ["+grp.Title()+"]")
		}
		calls = append(calls, f)
	}

	if args.HasModule("migration") {
		f, err = sql2.Migration(m, args, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL migration")
		}
		calls = append(calls, f)
	}
	if len(m.SeedData) > 0 {
		f, err = sql2.SeedData(m, args)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL seed data")
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
	return calls, nil
}

func basics(m *model2.Model, args *model2.Args, addHeader bool) (file.Files, error) {
	var calls file.Files
	f, err := gomodel2.Model(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render model")
	}
	calls = append(calls, f)

	if m.IsHistory() {
		f, err = gomodel2.History(m, args, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "can't render History")
		}
		calls = append(calls, f)
	}

	f, err = gomodel2.DTO(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render DTO")
	}
	calls = append(calls, f)

	fs, err := svc.ServiceAll(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render service")
	}
	calls = append(calls, fs...)

	f, err = controller2.Controller(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render controller")
	}
	calls = append(calls, f)
	return calls, nil
}
