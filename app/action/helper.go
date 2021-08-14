package action

import (
	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/file"
)

const (
	delimStart = "{{{"
	delimEnd   = "}}}"
)

func diffs(pm *PrjAndMods, addHeader bool) (file.Files, []*diff.Diff, error) {
	tgt := pm.PSvc.GetFilesystem(pm.Prj)

	srcFiles, err := pm.MSvc.GetFiles(pm.Mods, addHeader, tgt)
	if err != nil {
		return nil, nil, err
	}

	for _, fl := range srcFiles {
		err = file.ReplaceSections(fl, tgt)
		if err != nil {
			return nil, nil, err
		}
	}

	ctx := pm.Prj.ToTemplateContext()
	for _, f := range srcFiles {
		f.Content, err = runTemplate(f, ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	diffs := diff.FileLoader(srcFiles, tgt, false, pm.Logger)

	return srcFiles, diffs, nil
}
