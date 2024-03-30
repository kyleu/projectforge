package files

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files/controller"
	"projectforge.dev/projectforge/app/project/export/files/gomodel"
	"projectforge.dev/projectforge/app/project/export/files/sql"
	"projectforge.dev/projectforge/app/project/export/files/svc"
	"projectforge.dev/projectforge/app/project/export/files/view"
	"projectforge.dev/projectforge/app/project/export/model"
)

func ModelAll(m *model.Model, p *project.Project, args *model.Args, addHeader bool, linebreak string) (file.Files, error) {
	var calls file.Files
	var f *file.File

	fs, err := basics(m, args, addHeader, p.GoVersion(), linebreak)
	if err != nil {
		return nil, err
	}
	calls = append(calls, fs...)

	if args.HasModule("migration") && !m.HasTag("external") {
		f, err = sql.Migration(m, args, addHeader, linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL migration")
		}
		calls = append(calls, f)
	}
	if len(m.SeedData) > 0 {
		f, err = sql.SeedData(m, args.Database, linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render SQL seed data")
		}
		calls = append(calls, f)
	}

	fs, err = view.All(m, p, args, addHeader, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render list template")
	}
	calls = append(calls, fs...)

	return calls, nil
}

func basics(m *model.Model, args *model.Args, addHeader bool, goVersion string, linebreak string) (file.Files, error) {
	var calls file.Files
	f, err := gomodel.Model(m, args, addHeader, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render model")
	}
	fd, err := gomodel.ModelDiff(m, args, addHeader, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render model")
	}
	fm, err := gomodel.ModelMap(m, args, addHeader, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render model")
	}
	calls = append(calls, f, fd, fm)

	f, err = gomodel.Models(m, args, addHeader, goVersion, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render models")
	}
	calls = append(calls, f)

	f, err = gomodel.Row(m, args, addHeader, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render Row")
	}
	calls = append(calls, f)

	fs, err := svc.ServiceAll(m, args, addHeader, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render service")
	}
	calls = append(calls, fs...)

	f, err = controller.Controller(m, args, addHeader, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render controller")
	}
	calls = append(calls, f)
	return calls, nil
}
