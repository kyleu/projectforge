package files

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files/controller"
	"projectforge.dev/projectforge/app/project/export/files/gomodel"
	"projectforge.dev/projectforge/app/project/export/files/sql"
	"projectforge.dev/projectforge/app/project/export/files/svc"
	"projectforge.dev/projectforge/app/project/export/files/view"
)

func ModelAll(m *model.Model, p *project.Project, linebreak string) (file.Files, error) {
	args := p.ExportArgs
	var calls file.Files
	var f *file.File

	fs, err := basics(m, args, p.GoVersion(), linebreak)
	if err != nil {
		return nil, err
	}
	calls = append(calls, fs...)

	if !m.SkipDatabase() && !m.HasTag("external") {
		if args.HasModule("migration") {
			f, err = sql.Migration(m, args, linebreak)
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
	}

	if !m.SkipController() {
		fs, err = view.All(m, p, args, linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render list template")
		}
		calls = append(calls, fs...)
	}

	return calls, nil
}

func basics(m *model.Model, args *metamodel.Args, goVersion string, linebreak string) (file.Files, error) {
	var calls file.Files
	f, err := gomodel.Model(m, args, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render model")
	}
	fd, err := gomodel.ModelDiff(m, args, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render model")
	}
	fm, err := gomodel.ModelMap(m, args, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render model")
	}
	calls = append(calls, f, fd, fm)
	fa, err := gomodel.Models(m, args, goVersion, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render models")
	}
	calls = append(calls, fa)

	if !m.SkipDatabase() {
		f, err = gomodel.Row(m, args, linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render Row")
		}
		calls = append(calls, f)
	}

	if !m.SkipService() {
		fs, err := svc.ServiceAll(m, args, linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render service")
		}
		calls = append(calls, fs...)
	}

	if !m.SkipController() {
		f, err = controller.Controller(m, args, linebreak)
		if err != nil {
			return nil, errors.Wrap(err, "can't render controller")
		}
		calls = append(calls, f)
	}

	return calls, nil
}
