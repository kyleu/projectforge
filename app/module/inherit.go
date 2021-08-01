package module

import (
	"strings"

	"github.com/kyleu/projectforge/app/file"
	"github.com/pkg/errors"
)

func applyInheritance(fl *file.File, inh *file.Inheritiance, prior string) error {
	pIdx := 0
	if inh.Prefix != "" {
		pIdx = strings.Index(prior, inh.Prefix)
		if pIdx == -1 {
			return errors.Errorf("file [%s] does not contain prefix [%s]", fl.FullPath(), inh.Prefix)
		}
		pIdx += len(inh.Prefix)
	}

	sIdx := 0
	if inh.Suffix != "" {
		sIdx = strings.Index(prior[pIdx:], inh.Suffix)
		if sIdx == -1 {
			return errors.Errorf("file [%s] does not contain suffix [%s]", fl.FullPath(), inh.Suffix)
		}
		sIdx += pIdx
	}

	if pIdx == 0 && sIdx == 0 {
		return errors.Errorf("file [%s] must specify prefix or suffix", fl.FullPath())
	}

	fl.Content = prior[:pIdx] + inh.Content + prior[sIdx:]
	return nil
}
