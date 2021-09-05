package diff

import (
	"fmt"
	"strings"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/filesystem"
	"go.uber.org/zap"
)

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
		d := Calc(ret.Path, src.Content, tgt.Content)
		if len(d.Changes) > 0 {
			ret.Patch = d.Patch
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
			d = File(tgtFile, s)
		}
		if includeUnchanged || d.Status != StatusIdentical {
			ret = append(ret, d)
		}
	}
	return ret
}
