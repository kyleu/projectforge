package view

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project/export/model"
)

func All(m *model.Model, args *model.Args, addHeader bool) (file.Files, error) {
	var calls file.Files
	var f *file.File
	var err error

	f, err = list(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render list template")
	}
	calls = append(calls, f)

	f, err = table(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render table template")
	}
	calls = append(calls, f)

	f, err = detail(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render detail template")
	}
	calls = append(calls, f)

	f, err = edit(m, args, addHeader)
	if err != nil {
		return nil, errors.Wrap(err, "can't render edit template")
	}
	calls = append(calls, f)

	for _, grp := range m.GroupedColumns() {
		f, err = Grouping(m, grp, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "can't view controller for group ["+grp.Title()+"]")
		}
		calls = append(calls, f)
	}

	if m.IsHistory() {
		f, err = history(m, args, addHeader)
		if err != nil {
			return nil, errors.Wrap(err, "can't render history template")
		}
		calls = append(calls, f)
	}

	return calls, nil
}
