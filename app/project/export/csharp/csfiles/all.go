package csfiles

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

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

	modelGroups := lo.GroupBy(args.Models, func(x *model.Model) string {
		return strings.Join(x.Group, "/")
	})
	for grp, mdls := range modelGroups {
		if grp == "" {
			keys := lo.Map(mdls, func(x *model.Model, _ int) string {
				return x.Name
			})
			return nil, errors.Errorf("models [%s] have no group assigned, and can't be exported", strings.Join(keys, ", "))
		}
		if len(mdls) == 0 {
			continue
		}
		if strings.Contains(grp, "/") {
			return nil, errors.Errorf("invalid group [%s]", grp)
		}
		ns := CSNamespace(p.Key, grp)
		menuFile, err := menu(ns, mdls)
		if err != nil {
			return nil, err
		}
		ret = append(ret, menuFile)

		for _, m := range mdls {
			calls, err := ModelAll(ns, m, args)
			if err != nil {
				return nil, errors.Wrapf(err, "error processing model [%s]", m.Name)
			}
			ret = append(ret, calls...)
		}
	}

	return ret.Sort(), nil
}

func ModelAll(ns string, m *model.Model, args *model.Args) (file.Files, error) {
	ctrlr, err := controller(ns, m)
	if err != nil {
		return nil, err
	}
	svc, err := service(ns, m, args)
	if err != nil {
		return nil, err
	}
	list, err := cshtmlList(ns, m, args)
	if err != nil {
		return nil, err
	}
	detail, err := cshtmlDetail(ns, m, args)
	if err != nil {
		return nil, err
	}

	return file.Files{ctrlr, svc, list, detail}, nil
}
