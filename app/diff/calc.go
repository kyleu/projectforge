package diff

import (
	"fmt"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

type Result struct {
	Filename string    `json:"filename"`
	Src      string    `json:"src"`
	Tgt      string    `json:"tgt"`
	Diffs    Diffs     `json:"diffs"`
	Changes  []*Change `json:"changes"`
	Patch    string    `json:"patch"`
}

type Line struct {
	T string `json:"t"`
	V string `json:"v"`
}

type Change struct {
	From  int `json:"from"`
	To    int `json:"to"`
	Lines []*Line
}

func Calc(fn string, src string, tgt string) *Result {
	diffs := myers.ComputeEdits(span.URIFromPath(""), tgt, src)
	p, c := changes(tgt, diffs)
	return &Result{Filename: fn, Src: src, Tgt: tgt, Diffs: diffs, Changes: c, Patch: p}
}

func changes(src string, diffs []gotextdiff.TextEdit) (string, []*Change) {
	u := gotextdiff.ToUnified("", "", src, diffs)
	ret := make([]*Change, 0, len(u.Hunks))
	for _, h := range u.Hunks {
		lines := make([]*Line, 0, len(h.Lines))
		for _, l := range h.Lines {
			t := "unknown"
			switch l.Kind {
			case 0:
				t = "deleted"
			case 1:
				t = "added"
			case 2:
				t = "context"
			}
			lines = append(lines, &Line{T: t, V: l.Content})
		}
		ret = append(ret, &Change{From: h.FromLine, To: h.ToLine, Lines: lines})
	}
	patch := fmt.Sprint(u)
	patch = strings.TrimPrefix(patch, "--- ")
	patch = strings.TrimPrefix(patch, "\n")
	patch = strings.TrimPrefix(patch, "+++ ")
	patch = strings.TrimPrefix(patch, "\n")
	return patch, ret
}
