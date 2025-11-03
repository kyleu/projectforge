package diff

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/file"
	"projectforge.dev/projectforge/app/lib/filesystem"
	"projectforge.dev/projectforge/app/util"
)

func FileLoader(mods []string, src file.Files, tgt filesystem.FileLoader, ignoredFiles []string, includeUnchanged bool, logger util.Logger) (Diffs, error) {
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

			tgtFile = file.NewFile(p, t.Mode, b)
			if slices.Contains(ignoredFiles, tgtFile.FullPath()) {
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

		tgtMods := util.StringSplitAndTrim(line[open+1:cl], ",")

		hasAllMods := true
		lo.ForEach(tgtMods, func(mod string, _ int) {
			if !lo.Contains(mods, mod) {
				hasAllMods = false
			}
		})
		if hasAllMods {
			if tgtFile != nil {
				tgtFile.Content = util.StringJoin(slices.Delete(lines, lineIdx, lineIdx), util.StringDetectLinebreak(s.Content))
			}
		} else {
			return false, nil
		}
	}
	return true, nil
}
