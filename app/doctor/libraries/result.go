package libraries

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/google/uuid"

	"projectforge.dev/projectforge/app/util"
)

type Result struct {
	ID      uuid.UUID `json:"id,omitzero"`
	Library *Library  `json:"library,omitzero"`
	Action  string    `json:"action,omitzero"`
	Output  []any     `json:"output,omitzero"`
	Outcome string    `json:"outcome,omitzero"`
}

func NewResult(library *Library, action string) *Result {
	return &Result{ID: util.UUID(), Library: library, Action: action}
}

func (r *Result) AddMessage(s string, args ...any) {
	r.Output = append(r.Output, fmt.Sprintf(s, args...))
}

func (r *Result) String() string {
	if r.Outcome == "" {
		return fmt.Sprintf("%s: No Result", r.Library.String())
	}
	return fmt.Sprintf("%s: %s", r.Library.String(), r.Outcome)
}

type Results []*Result

func (r Results) Sort() {
	slices.SortFunc(r, func(l *Result, r *Result) int {
		li, ri := AllLibraries.Index(l.Library.Key), AllLibraries.Index(r.Library.Key)
		return cmp.Compare(li, ri)
	})
}
