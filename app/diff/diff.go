package diff

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/sergi/go-diff/diffmatchpatch"
	"go.uber.org/zap"
)

var dmp = diffmatchpatch.New()

type Diff struct {
	Path   string  `json:"path"`
	Status *Status `json:"status"`
	Patch  string  `json:"patch,omitempty"`
}

func (d *Diff) String() string {
	return fmt.Sprintf("%s:%s", d.Path, d.Status)
}

func File(src *file.File, tgt *file.File) *Diff {
	ret := &Diff{Path: src.FullPath()}
	switch {
	case src == nil:
		ret.Status = StatusMissing
	case tgt == nil:
		ret.Status = StatusNew
	default:
		t1, t2, a := dmp.DiffLinesToChars(tgt.Content, src.Content)
		diffs := dmp.DiffMain(t1, t2, true)
		diffs = dmp.DiffCharsToLines(diffs, a)
		diffs = dmp.DiffCleanupSemantic(diffs)
		patches := dmp.PatchMake(tgt.Content, diffs)
		if len(patches) > 0 {
			ret.Patch = dmp.PatchToText(patches)
			pte, err := url.QueryUnescape(ret.Patch)
			if err == nil {
				ret.Patch = pte
			}
			ret.Status = StatusDifferent
		} else {
			ret.Status = StatusIdentical
		}
	}
	return ret
}

func Files(src file.Files, tgt file.Files, includeUnchanged bool) []*Diff {
	var ret []*Diff
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

func FileLoader(src file.Files, tgt filesystem.FileLoader, includeUnchanged bool, logger *zap.SugaredLogger) []*Diff {
	var ret []*Diff
	for _, s := range src {
		p := s.FullPath()
		t, _ := tgt.Stat(p)

		skip := false
		var tgtFile *file.File
		if t != nil {
			b, err := tgt.ReadFile(p)
			if err != nil {
				msg := "An error was encountered: %+v"
				ret = append(ret, &Diff{Path: p, Status: &Status{Key: "error", Title: fmt.Sprintf(msg, err)}})
			}

			tgtFile = file.NewFile(p, t.Mode(), b, false, logger)

			linefeed := strings.Index(tgtFile.Content, "\n")
			if linefeed > -1 {
				firstLine := tgtFile.Content[:linefeed]
				if strings.Contains(firstLine, "$PF_IGNORE$") {
					skip = true
				}
			}
		}
		var d *Diff
		if skip {
			d = &Diff{Path: s.FullPath(), Status: StatusSkipped}
		} else {
			d = File(s, tgtFile)
		}
		if includeUnchanged || d.Status != StatusIdentical {
			ret = append(ret, d)
		}
	}
	return ret
}
