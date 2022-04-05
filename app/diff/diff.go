package diff

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	var ret Diffs
	for _, s := range src {
		p := s.FullPath()
		t := tgt.Get(p)
		d := File(s, t)
		if includeUnchanged || d.Status != StatusIdentical {
			ret = append(ret, d)
		}
	}
	return ret
}

func FileLoader(mods []string, src file.Files, tgt filesystem.FileLoader, includeUnchanged bool, logger *zap.SugaredLogger) (Diffs, error) {
	var ret Diffs
	for _, s := range src {
		p := s.FullPath()
		t, _ := tgt.Stat(p)

		skip := false
		var tgtFile *file.File
		if t != nil {
			b, err := tgt.ReadFile(p)
			if err != nil {
				ret = append(ret, &Diff{Path: p, Status: &Status{Key: "error", Title: fmt.Sprintf("An error was encountered: %+v", err)}})
			}

			tgtFile = file.NewFile(p, t.Mode(), b, false, logger)

			if strings.Contains(tgtFile.Content, file.IgnorePattern) {
				skip = true
			}
		}

		if idx := strings.Index(s.Content, file.ModulePrefix); idx > 1 {
			lines := strings.Split(s.Content, "\n")
			lineIdx := -1
			for idx, line := range lines {
				if strings.Contains(line, file.ModulePrefix) {
					lineIdx = idx
					break
				}
			}
			line := lines[lineIdx]
			if !strings.Contains(line, file.ModulePrefix) {
				return nil, errors.New("module requirement tag must be on first meaningful line of the file")
			}

			open, cl := strings.Index(line, "("), strings.Index(line, ")")
			if open == -1 || cl == -1 {
				return nil, errors.New("module requirement tag must contain parentheses")
			}
			var hasAllMods = true
			for _, mod := range util.StringSplitAndTrim(line[open+1:cl], ",") {
				if !slices.Contains(mods, mod) {
					hasAllMods = false
				}
			}
			if hasAllMods {
				if tgtFile != nil {
					tgtFile.Content = strings.Join(slices.Delete(lines, lineIdx, lineIdx), "\n")
				}
			} else {
				skip = true
			}
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
