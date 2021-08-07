package file

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	inheritPrefix = prefix + "INHERIT$"
	prefixPattern = prefix + "PREFIX$"
	suffixPattern = prefix + "SUFFIX$"
)

type Inheritiance struct {
	Prefix  string `json:"prefix"`
	Suffix  string `json:"suffix"`
	Content string `json:"content"`
}

func InheritanceContent(fl *File) (*Inheritiance, error) {
	if !strings.Contains(fl.Content, inheritPrefix) {
		return nil, nil
	}
	lines := strings.Split(fl.Content, "\n")
	ret := &Inheritiance{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		hIdx := strings.Index(line, headerContent)
		if hIdx > -1 {
			continue
		}
		inhIdx := strings.Index(line, inheritPrefix)
		if inhIdx > -1 {
			continue
		}
		prefixIdx := strings.Index(line, prefixPattern)
		if prefixIdx > -1 {
			ret.Prefix = line[prefixIdx+len(prefixPattern):]
			ret.Prefix = strings.TrimSpace(ret.Prefix)
			ret.Prefix = strings.TrimPrefix(strings.TrimSuffix(ret.Prefix, `]`), `[`)
			continue
		}
		suffixIdx := strings.Index(line, suffixPattern)
		if suffixIdx > -1 {
			ret.Suffix = line[suffixIdx+len(suffixPattern):]
			ret.Suffix = strings.TrimSpace(ret.Suffix)
			ret.Suffix = strings.TrimPrefix(strings.TrimSuffix(ret.Suffix, `]`), `[`)
			continue
		}
		ret.Content += "\n" + line
	}
	if ret.Prefix == "" && ret.Suffix == "" {
		return nil, errors.Errorf("must provide one of [%s, %s] to enable inheritance", ret.Prefix, ret.Suffix)
	}
	ret.Content += "\n"
	return ret, nil
}
