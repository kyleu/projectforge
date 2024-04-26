package csfiles

import (
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/model"
)

func CSAll(p *project.Project) (file.Files, error) {
	if p.ExportArgs == nil {
		return nil, errors.New("export arguments aren't loaded")
	}
	args := p.ExportArgs
	if err := args.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid export arguments")
	}
	var ret file.Files

	menuFile, err := menu(args.Models, p)
	if err != nil {
		return nil, err
	}
	ret = append(ret, menuFile)

	for _, m := range args.Models {
		calls, err := ModelAll(m, p)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing model [%s]", m.Name)
		}
		ret = append(ret, calls...)
	}

	return ret, nil
}

func ModelAll(m *model.Model, p *project.Project) (file.Files, error) {
	var ret file.Files

	ctrlr, err := controller(m, p)
	if err != nil {
		return nil, err
	}
	ret = append(ret, ctrlr)

	svc, err := service(m, p)
	if err != nil {
		return nil, err
	}
	ret = append(ret, svc)

	list, err := cshtmlList(m)
	if err != nil {
		return nil, err
	}
	ret = append(ret, list)

	return ret, nil
}
