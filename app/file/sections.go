package file

import (
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/lib/filesystem"
)

const (
	InjectPrefix  = prefix + "INJECT"
	sectionPrefix = prefix + "SECTION"
	startPattern  = "_START("
	endPattern    = "_END("
	closePattern  = ")$"
	IgnorePattern = prefix + "IGNORE" + suffix
)

func Inject(fl *File, content map[string]string) error {
	injectIndexes, err := sectionIndexes(fl.Content, InjectPrefix)
	if err != nil {
		return errors.Wrap(err, "can't read inject indexes")
	}

	if len(injectIndexes) < len(content) {
		// return errors.Errorf("file contains [%d] injection points, but [%d] are required", len(injectIndexes), len(content))
	}

	for _, idx := range injectIndexes {
		c, ok := content[idx.Key]
		if !ok {
			return errors.Errorf("no inject index for key [%s]", idx.Key)
		}
		fl.Content = fl.Content[0:idx.Start] + c + fl.Content[idx.End:]
	}
	return nil
}

func ReplaceSections(fl *File, tgt filesystem.FileLoader) error {
	f := fl.FullPath()
	if strings.Contains(fl.Content, sectionPrefix) && tgt.Exists(f) {
		tgtBytes, _ := tgt.ReadFile(f)
		if utf8.Valid(tgtBytes) {
			newContent, err := copySections(string(tgtBytes), fl.Content)
			if err != nil {
				return errors.Wrapf(err, "error reading sections from [%s]", f)
			}
			fl.Content = newContent
		}
	}
	return nil
}

func copySections(src string, tgt string) (string, error) {
	srcSections, err := sectionIndexes(src, sectionPrefix)
	if err != nil {
		return "", errors.Wrap(err, "unable to read sections from source content")
	}
	tgtSections, err := sectionIndexes(tgt, sectionPrefix)
	if err != nil {
		return "", errors.Wrap(err, "unable to read sections from target content")
	}

	ret := tgt
	for _, sec := range srcSections {
		tgtSec := tgtSections.Get(sec.Key)
		if tgtSec == nil {
			return "", errors.Errorf("no section [%s] found in target", sec.Key)
		}
		content := src[sec.Start:sec.End]
		ret = ret[0:tgtSec.Start] + content + ret[tgtSec.End:]
	}

	return ret, nil
}

type section struct {
	Key   string
	Start int
	End   int
}

type sections []*section

func (s sections) Get(key string) *section {
	for _, x := range s {
		if x.Key == key {
			return x
		}
	}
	return nil
}

func (s sections) Sort() {
	sort.Slice(s, func(i int, j int) bool {
		return s[i].Start > s[j].Start
	})
}

func sectionIndexes(s string, prefix string) (sections, error) {
	var ret sections

	for idx, c := range s {
		if c == '$' && len(s) > idx+len(prefix) {
			candidate := s[idx : idx+len(prefix)]
			if candidate == prefix {
				nextDollar := strings.Index(s[idx+len(prefix):], closePattern)
				if nextDollar == -1 {
					return nil, errors.New("found section, but no closing delimiter")
				}
				if nextDollar > 256 {
					return nil, errors.Errorf("found section, but no closing delimiter within [%d] bytes", nextDollar)
				}
				endIdx := idx + len(prefix) + nextDollar + len(closePattern)
				text := s[idx+len(prefix) : endIdx]
				var err error
				ret, err = parseText(ret, endIdx, text, prefix)
				if err != nil {
					return nil, errors.Wrapf(err, "unable to parse text [%s]", text)
				}
			}
		}
	}
	ret.Sort()
	return ret, nil
}

func parseText(ret sections, idx int, text string, prefix string) (sections, error) {
	switch {
	case strings.HasPrefix(text, startPattern):
		currSection := text[len(startPattern) : len(text)-len(closePattern)]
		if ret.Get(currSection) != nil {
			return nil, errors.Errorf("multiple sections found with key [%s]", currSection)
		}
		ret = append(ret, &section{Key: currSection, Start: idx})
		return ret, nil
	case strings.HasPrefix(text, endPattern):
		if len(ret) == 0 {
			return nil, errors.New("encountered end section pattern before start")
		}
		curr := ret[len(ret)-1]
		sec := text[len(endPattern) : len(text)-len(closePattern)]
		if curr.Key != sec {
			return nil, errors.Errorf("encountered nested section patterns (%s within %s)", sec, curr.Key)
		}
		curr.End = idx - len(prefix) - len(text)
		return ret, nil
	default:
		return nil, errors.Errorf("invalid section pattern [%s]", text)
	}
}
