package action

import (
	"path"
	"strings"

	"github.com/kyleu/projectforge/app/diff"
	"github.com/kyleu/projectforge/app/file"
)

const (
	delimStart = "{{{"
	delimEnd   = "}}}"
)

func diffs(pm *PrjAndMods) (file.Files, diff.Diffs, error) {
	tgt := pm.PSvc.GetFilesystem(pm.Prj)

	srcFiles, err := pm.MSvc.GetFiles(pm.Mods)
	if err != nil {
		return nil, nil, err
	}

	portOffsets := map[string]int{}
	for _, m := range pm.Prj.Modules {
		for k, v := range pm.Mods.Get(m).PortOffsets {
			portOffsets[k] = v
		}
	}

	ctx := pm.Prj.ToTemplateContext(portOffsets)

	for _, f := range srcFiles {
		origPath := f.FullPath()
		if strings.Contains(origPath, delimStart) {
			newPath, e := runTemplate("filename", origPath, ctx)
			if e != nil {
				return nil, nil, e
			}
			p, n := path.Split(newPath)
			f.Path = strings.Split(p, "/")
			f.Name = n
		}
		err = file.ReplaceSections(f, tgt)
		if err != nil {
			return nil, nil, err
		}
	}

	for _, f := range srcFiles {
		f.Content, err = runTemplateFile(f, ctx)
		if err != nil {
			return nil, nil, err
		}
	}

	diffs := diff.FileLoader(srcFiles, tgt, false, pm.Logger)

	return srcFiles, diffs, nil
}
