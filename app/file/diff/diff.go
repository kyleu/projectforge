package diff

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

type Diff struct {
	Path    string  `json:"path"`
	Status  *Status `json:"status"`
	Patch   string  `json:"patch,omitempty"`
	Changes Changes `json:"changes,omitempty"`
}

func (d *Diff) String() string {
	return fmt.Sprintf("%s:%s", d.Path, d.Status)
}

type Diffs []*Diff

func (d Diffs) HasStatus(s *Status) bool {
	return lo.ContainsBy(d, func(x *Diff) bool {
		return x.Status.Key == s.Key
	})
}

func (d Diffs) Paths() []string {
	return lo.Uniq(lo.Map(d, func(x *Diff, _ int) string {
		return x.Path
	}))
}

func File(src *file.File, tgt *file.File) *Diff {
	ret := &Diff{Path: src.FullPath()}
	switch {
	case src == nil:
		ret.Status = StatusMissing
	case tgt == nil:
		ret.Status = StatusNew
	default:
		d := Calc(ret.Path, src.Content, tgt.Content)
		if len(d.Changes) > 0 {
			ret.Patch = d.Patch
			ret.Status = StatusDifferent
			ret.Changes = d.Changes
		} else {
			ret.Status = StatusIdentical
		}
	}
	return ret
}

func Files(src file.Files, tgt file.Files, includeUnchanged bool) Diffs {
	return lo.FlatMap(src, func(s *file.File, _ int) []*Diff {
		p := s.FullPath()
		t := tgt.Get(p)
		d := File(s, t)
		if includeUnchanged || d.Status != StatusIdentical {
			return Diffs{d}
		}
		return nil
	})
}

func FileLoader(mods []string, src file.Files, tgt filesystem.FileLoader, includeUnchanged bool, logger util.Logger) (Diffs, error) {
	var ret Diffs
	for _, s := range src {
		p := s.FullPath()
		t, _ := tgt.Stat(p)

		skip := false
		var tgtFile *file.File
		if t != nil {
			b, err := tgt.ReadFile(p)
			if err != nil {
				ret = append(ret, &Diff{Path: p, Status: &Status{Key: util.KeyError, Title: fmt.Sprintf("An error was encountered: %+v", err)}})
			}

			tgtFile = file.NewFile(p, t.Mode(), b, false, logger)

			if strings.Contains(tgtFile.Content, file.IgnorePattern) {
				skip = true
			}
			if strings.Contains(s.Content, file.GenerateOncePattern) {
				skip = true
			}
		}

		matches, err := matchesModules(s, mods, tgtFile)
		if err != nil {
			return nil, err
		}
		if !matches {
			skip = true
		}

		var d *Diff
		if skip {
			d = &Diff{Path: s.FullPath(), Status: StatusSkipped}
		} else {
			d = File(s, tgtFile)
		}
		if includeUnchanged || (d.Status != StatusIdentical && d.Status != StatusSkipped) {
			ret = append(ret, d)
		}
	}
	return ret, nil
}

func matchesModules(s *file.File, mods []string, tgtFile *file.File) (bool, error) {
	if idx := strings.Index(s.Content, file.ModulePrefix); idx > 1 {
		lines := util.StringSplitLines(s.Content)
		line, lineIdx, _ := lo.FindIndexOf(lines, func(line string) bool {
			return strings.Contains(line, file.ModulePrefix)
		})
		if !strings.Contains(line, file.ModulePrefix) {
			return false, errors.New("module requirement tag must be on first meaningful line of the file")
		}

		open, cl := strings.Index(line, "("), strings.Index(line, ")")
		if open == -1 || cl == -1 {
			return false, errors.New("module requirement tag must contain parentheses")
		}
		hasAllMods := true
		lo.ForEach(util.StringSplitAndTrim(line[open+1:cl], ","), func(mod string, _ int) {
			if !lo.Contains(mods, mod) {
				hasAllMods = false
			}
		})
		if hasAllMods {
			if tgtFile != nil {
				tgtFile.Content = strings.Join(slices.Delete(lines, lineIdx, lineIdx), util.StringDetectLinebreak(s.Content))
			}
		} else {
			return false, nil
		}
	}
	return true, nil
}
