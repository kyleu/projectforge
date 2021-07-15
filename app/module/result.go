package module

import (
	"github.com/kyleu/projectforge/app/file"
)

type Result struct {
	Key      string       `json:"key"`
	Status   string       `json:"status"`
	Diffs    []*file.Diff `json:"diffs,omitempty"`
	Duration int          `json:"duration,omitempty"`
}

type Results []*Result
