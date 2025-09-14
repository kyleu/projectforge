package files

import (
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/files/goevent"
)

func EventAll(p *project.Project, events model.Events, linebreak string) (file.Files, error) {
	var ret file.Files
	for _, evt := range events {
		f, err := goevent.Event(evt, p.ExportArgs, linebreak)
		if err != nil {
			return nil, errors.Wrapf(err, "error processing event [%s]", evt.Name)
		}
		ret = append(ret, f)
	}
	return ret, nil
}
