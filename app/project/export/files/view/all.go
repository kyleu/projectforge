package view

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project"
)

func All(m *model.Model, p *project.Project, args *metamodel.Args, linebreak string) (file.Files, error) {
	var calls file.Files
	var f *file.File
	var err error

	f, err = list(m, args, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render list template")
	}
	calls = append(calls, f)

	f, err = table(m, args, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render table template")
	}
	calls = append(calls, f)

	f, err = detail(m, args, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render detail template")
	}
	calls = append(calls, f)

	f, err = edit(m, p, args, linebreak)
	if err != nil {
		return nil, errors.Wrap(err, "can't render edit template")
	}
	calls = append(calls, f)

	return calls, nil
}
