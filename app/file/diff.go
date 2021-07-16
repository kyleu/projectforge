package file

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/sergi/go-diff/diffmatchpatch"
)

var dmp = diffmatchpatch.New()

const (
	StatusDifferent = "different"
	StatusIdentical = "identical"
	StatusMissing   = "missing"
	StatusNew       = "new"
	StatusSkipped   = "skipped"
)

type Diff struct {
	Path   string `json:"path"`
	Status string `json:"status"`
	Patch  string `json:"patch,omitempty"`
}

func (d *Diff) String() string {
	return fmt.Sprintf("%s:%s", d.Path, d.Status)
}

func (f *File) Diff(tgt *File) *Diff {
	ret := &Diff{Path: f.FullPath()}
	switch {
	case f == nil:
		ret.Status = StatusMissing
	case tgt == nil:
		ret.Status = StatusNew
	default:
		t1, t2, a := dmp.DiffLinesToChars(tgt.Content, f.Content)
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

func (f Files) Diff(tgt Files, includeUnchanged bool) []*Diff {
	var ret []*Diff
	for _, file := range f {
		p := file.FullPath()
		t := tgt.Get(p)
		d := file.Diff(t)
		if includeUnchanged || d.Status != StatusIdentical {
			ret = append(ret, d)
		}
	}
	return ret
}

func (f Files) DiffFileLoader(tgt filesystem.FileLoader, includeUnchanged bool) []*Diff {
	var ret []*Diff
	for _, file := range f {
		p := file.FullPath()
		t := tgt.Stat(p)

		skip := false
		var tgtFile *File

		if t != nil {
			b, err := tgt.ReadFile(p)
			if err != nil {
				ret = append(ret, &Diff{Path: p, Status: fmt.Sprintf("error: %+v", err)})
			}

			content := string(b)
			linefeed := strings.Index(content, "\n")
			if linefeed > -1 {
				firstLine := content[:linefeed]
				if strings.Contains(firstLine, "$PF_IGNORE$") {
					skip = true
				}
			}
			tgtFile = NewFile(p, t.Mode(), b)
		}
		var d *Diff
		if skip {
			d = &Diff{Path: file.FullPath(), Status: StatusSkipped}
		} else {
			d = file.Diff(tgtFile)
		}
		if includeUnchanged || d.Status != StatusIdentical {
			ret = append(ret, d)
		}
	}
	return ret
}
