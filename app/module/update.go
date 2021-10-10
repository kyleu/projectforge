package module

import (
	"github.com/kyleu/projectforge/app/diff"
)

func (s *Service) UpdateFile(mods Modules, d *diff.Diff) (error) {
	for _, mod := range mods {
		loader := s.GetFilesystem(mod.Key)
		fs, err := loader.ListFilesRecursive("", nil)
		if err != nil {
			return err
		}
		for _, f := range fs {
			if f == configFilename {
				continue
			}
			if f != d.Path {
				continue
			}
			mode, b, err := fileContent(loader, f)
			if err != nil {
				return err
			}

			newContent, err := applyDiff(b, d)
			if len(newContent) > 0 {
				err = loader.WriteFile(d.Path, newContent, mode, true)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func applyDiff(b []byte, d *diff.Diff) ([]byte, error) {
	return nil, nil
}
