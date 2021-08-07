package file

import (
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/pkg/errors"
)

const (
	sectionPrefix = prefix + "SECTION"
	startPattern  = "_START("
	endPattern    = "_END("
	closePattern  = ")$"
)

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
	srcSections, err := sectionIndexes(src)
	if err != nil {
		return "", errors.Wrap(err, "unable to read sections from source content")
	}
	tgtSections, err := sectionIndexes(tgt)
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
	sort.Slice(s, func(i, j int) bool {
		return s[i].Start > s[j].Start
	})
}

func sectionIndexes(s string) (sections, error) {
	var ret sections
	parseText := func(idx int, text string) error {
		switch {
		case strings.HasPrefix(text, startPattern):
			currSection := text[len(startPattern) : len(text)-len(closePattern)]
			if ret.Get(currSection) != nil {
				return errors.Errorf("multiple sections found with key [%s]", currSection)
			}
			ret = append(ret, &section{Key: currSection, Start: idx})
		case strings.HasPrefix(text, endPattern):
			if len(ret) == 0 {
				return errors.New("encountered end section pattern before start")
			}
			curr := ret[len(ret)-1]
			sec := text[len(endPattern) : len(text)-len(closePattern)]
			if curr.Key != sec {
				return errors.Errorf("encountered nested section patterns (%s within %s)", sec, curr.Key)
			}
			curr.End = idx - len(sectionPrefix) - len(text)
		default:
			return errors.Errorf("invalid section pattern [%s]", text)
		}
		return nil
	}

	for idx, c := range s {
		if c == '$' && len(s) > idx+len(sectionPrefix) {
			candidate := s[idx : idx+len(sectionPrefix)]
			if candidate == sectionPrefix {
				nextDollar := strings.Index(s[idx+len(sectionPrefix):], closePattern)
				if nextDollar == -1 {
					return nil, errors.New("found section, but no closing delimiter")
				}
				if nextDollar > 256 {
					return nil, errors.Errorf("found section, but no closing delimiter within [%d] bytes", nextDollar)
				}
				endIdx := idx + len(sectionPrefix) + nextDollar + len(closePattern)
				text := s[idx+len(sectionPrefix) : endIdx]
				err := parseText(endIdx, text)
				if err != nil {
					return nil, errors.Wrapf(err, "unable to parse text [%s]", text)
				}
			}
		}
	}
	ret.Sort()
	return ret, nil
}
