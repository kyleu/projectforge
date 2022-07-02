package diff

import (
	"github.com/hexops/gotextdiff"
)

type Edits []gotextdiff.TextEdit

type Example struct {
	File     string `json:"title"`
	Src      string `json:"src"`
	Tgt      string `json:"tgt"`
	Expected Edits  `json:"expected"`
}

func NewExample(fn string, src string, tgt string) *Example {
	return &Example{File: fn, Src: src, Tgt: tgt}
}

func (e *Example) Calc() *Result {
	return Calc(e.File, e.Src, e.Tgt)
}

const src = "header\nbody\nfooter\n"

var (
	matching     = NewExample("matching", src, src)
	added        = NewExample("added", src, "header\nbody\nadded\nfooter\n")
	removed      = NewExample("removed", src, "header\nfooter\n")
	startChange  = NewExample("startChange", src, "hdr\nbody\nfooter\n")
	middleChange = NewExample("middleChange", src, "header\nb\nfooter\n")
	endChange    = NewExample("endChange", src, "header\nbody\nft\n")
	noIntersect  = NewExample("noIntersect", src, "x\ny\nz\n")
)

var AllExamples = []*Example{matching, added, removed, startChange, middleChange, endChange, noIntersect}
