package files

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files/controller"
	"projectforge.dev/projectforge/app/project/export/files/gomodel"
	"projectforge.dev/projectforge/app/project/export/files/grpc"
	"projectforge.dev/projectforge/app/project/export/files/sql"
	"projectforge.dev/projectforge/app/project/export/files/svc"
	"projectforge.dev/projectforge/app/project/export/files/view"
	"projectforge.dev/projectforge/app/project/export/model"
)

func ModelAll(m *model.Model, p *project.Project, args *model.Args, addHeader bool) (file.Files, error) {
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
	if len(m.SeedData) > 0 {
		f, err = sql.SeedData(m, args)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL seed data")
		}
		calls = append(calls, f)
	}

	fs, err = view.All(m, p, args, addHeader)
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

func basics(m *model.Model, args *model.Args, addHeader bool) (file.Files, error) {
	var calls file.Files
	f, err := gomodel.Model(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render model")
	}
	calls = append(calls, f)

	f, err = gomodel.Models(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render models")
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
