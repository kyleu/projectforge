package action

import (
	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/module"
	"github.com/kyleu/projectforge/app/project"
	"go.uber.org/zap"
)

const (
	delimStart = "{{{"
	delimEnd   = "}}}"
)

func diffs(prj *project.Project, mods module.Modules, addHeader bool, mSvc *module.Service, pSvc *project.Service, logger *zap.SugaredLogger) (file.Files, []*diff.Diff, error) {
	tgt := pSvc.GetFilesystem(prj)

	srcFiles, err := mSvc.GetFiles(mods, addHeader, tgt)
	if err != nil {
		return nil, nil, err
	}

	for _, fl := range srcFiles {
		err = file.ReplaceSections(fl, tgt)
		if err != nil {
			return nil, nil, err
		}
	}

	ctx := prj.ToTemplateContext()
	for _, f := range srcFiles {
		f.Content, err = runTemplate(f, ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	diffs := diff.FileLoader(srcFiles, tgt, false, logger)

	return srcFiles, diffs, nil
}
