package diff

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

const (
	addedKey   = "added"
	deletedKey = "deleted"
	contextKey = "context"
)

type Result struct {
	Filename string  `json:"filename"`
	Src      string  `json:"src"`
	Tgt      string  `json:"tgt"`
	Edits    Edits   `json:"edits"`
	Changes  Changes `json:"changes"`
	Patch    string  `json:"patch"`
}

type Results []*Result

type Line struct {
	T string `json:"t"`
	V string `json:"v"`
}

type Lines []*Line

func (l Line) String() string {
	switch l.T {
	case addedKey:
		return " + " + l.V
	case deletedKey:
		return " - " + l.V
	case contextKey:
		return " . " + l.V
	default:
		return " ? " + l.V
	}
}

type Change struct {
	From  int `json:"from"`
	To    int `json:"to"`
	Lines Lines
}

type Changes []*Change

func Calc(fn string, src string, tgt string) *Result {
	edits := myers.ComputeEdits(span.URIFromPath(""), tgt, src)
	p, c := changes(tgt, edits)
	return &Result{Filename: fn, Src: src, Tgt: tgt, Edits: edits, Changes: c, Patch: p}
}

func changes(src string, edits []gotextdiff.TextEdit) (string, Changes) {
	u := gotextdiff.ToUnified("", "", src, edits)
	ret := make(Changes, 0, len(u.Hunks))
	lo.ForEach(u.Hunks, func(h *gotextdiff.Hunk, _ int) {
		lines := lo.Map(h.Lines, func(l gotextdiff.Line, _ int) *Line {
			t := "unknown"
			switch l.Kind {
			case gotextdiff.Delete:
				t = deletedKey
			case gotextdiff.Insert:
				t = addedKey
			case gotextdiff.Equal:
				t = contextKey
			}
			return &Line{T: t, V: l.Content}
		})
		ret = append(ret, &Change{From: h.FromLine, To: h.ToLine, Lines: lines})
	})
	patch := fmt.Sprint(u)
	patch = strings.TrimPrefix(patch, "--- ")
	patch = strings.TrimPrefix(patch, "\n")
	patch = strings.TrimPrefix(patch, "+++ ")
	patch = strings.TrimPrefix(patch, "\n")
	return patch, ret
}
