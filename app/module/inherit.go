package module

import (
	"strings"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func applyInheritance(fl *file.File, inh *file.Inheritiance, addHeader bool, loader filesystem.FileLoader, logger *zap.SugaredLogger) error {
	cand, prior, err := getPrior(fl, loader)
	if err != nil {
		return err
	}
	newFile := file.NewFile(fl.FullPath(), fl.Mode, prior, addHeader, logger)
	nc := newFile.Content

	newInh, err := file.InheritanceContent(fl)
	if err != nil {
		return err
	}

	if newInh != nil {
		return applyInheritance(newFile, newInh, addHeader, cand, logger)
	}

	if prior == nil {
		return errors.Errorf("missing nested file [%s] for inheritance in [%s]", fl.FullPath(), loader.Root())
	}

	pIdx := 0
	if inh.Prefix != "" {
		pIdx = strings.Index(nc, inh.Prefix)
		if pIdx == -1 {
			return errors.Errorf("file [%s] does not contain prefix [%s]", fl.FullPath(), inh.Prefix)
		}
		pIdx += len(inh.Prefix)
	}

	sIdx := 0
	if inh.Suffix != "" {
		sIdx = strings.Index(nc, inh.Suffix)
		if sIdx == -1 {
			return errors.Errorf("file [%s] does not contain suffix [%s]", fl.FullPath(), inh.Suffix)
		}
	}

	if pIdx == 0 && sIdx == 0 {
		return errors.Errorf("file [%s] must specify prefix or suffix", fl.FullPath())
	}

	fl.Content = nc[:pIdx] + inh.Content + nc[sIdx:]
	return nil
}

func getPrior(fl *file.File, loader filesystem.FileLoader) (filesystem.FileLoader, []byte, error) {
	var candidate filesystem.FileLoader
	var prior []byte
	for _, fs := range loader.GetChildren() {
		if fs.Exists(fl.FullPath()) {
			if prior != nil {
				return nil, nil, errors.Errorf("multiple nested inheritance files [%s] in [%s]", fl.FullPath(), loader.Root())
			}
			var err error
			prior, err = fs.ReadFile(fl.FullPath())
			if err != nil {
				return nil, nil, errors.Wrapf(err, "unable to read nested file [%s] in [%s]", fl.FullPath(), loader.Root())
			}
			candidate = fs
		}
	}
	return candidate, prior, nil
}

