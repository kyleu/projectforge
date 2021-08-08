package module

import (
	"github.com/kyleu/projectforge/app/diff"
)

type Result struct {
	Keys     []string     `json:"keys"`
	Status   string       `json:"status"`
	Diffs    []*diff.Diff `json:"diffs,omitempty"`
	Actions  Resolutions  `json:"actions,omitempty"`
	Duration int          `json:"duration,omitempty"`
}

type Results []*Result
