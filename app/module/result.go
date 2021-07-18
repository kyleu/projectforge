package module

import (
	"github.com/kyleu/projectforge/app/diff"
)

type Result struct {
	Key      string       `json:"key"`
	Status   string       `json:"status"`
	Diffs    []*diff.Diff `json:"diffs,omitempty"`
	Duration int          `json:"duration,omitempty"`
}

type Results []*Result
