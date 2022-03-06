package module

import (
	"bytes"
	"fmt"

	"projectforge.dev/app/diff"
)

func (s *Service) UpdateFile(mods Modules, d *diff.Diff) ([]string, error) {
	var ret []string
	for _, mod := range mods {
		loader := s.GetFilesystem(mod.Key)
		if !loader.Exists(d.Path) {
			continue
		}
		mode, b, err := fileContent(loader, d.Path)
		if err != nil {
			return nil, err
		}

		newContent, err := diff.ApplyInverse(b, d)
		if err != nil {
			return nil, err
		}
		if len(newContent) > 0 {
			if bytes.Equal(b, newContent) {
				ret = append(ret, fmt.Sprintf("no changes required to [%s] for module [%s]", d.Path, mod.Key))
			} else {
				ret = append(ret, fmt.Sprintf("wrote [%d] bytes to [%s] for module [%s]", len(newContent), d.Path, mod.Key))
				err = loader.WriteFile(d.Path, newContent, mode, true)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return ret, nil
}
