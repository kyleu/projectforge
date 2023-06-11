package git

import (
	"projectforge.dev/projectforge/app/util"
	"time"
)

type HistoryResult struct {
	Path    string         `json:"path,omitempty"`
	Since   *time.Time     `json:"since,omitempty"`
	Authors []string       `json:"authors,omitempty"`
	Limit   int            `json:"limit,omitempty"`
	Entries HistoryEntries `json:"entries"`
	Debug   any            `json:"debug,omitempty"`
}

type HistoryEntry struct {
	Headers     map[string]string `json:"headers" xml:"headers"`
	SHA         string            `json:"sha" xml:"sha"`
	AuthorName  string            `json:"authorName" xml:"authorName"`
	AuthorEmail string            `json:"authorEmail" xml:"authorEmail"`
	Message     string            `json:"message" xml:"message"`
	Occurred    string            `json:"occurred" xml:"occurred"`
	Files       HistoryFiles      `json:"files" xml:"files"`
}

func (h *HistoryEntry) OccurredTime() *time.Time {
	ret, _ := util.TimeFromVerbose(h.Occurred)
	return ret
}

type HistoryEntries []*HistoryEntry

type HistoryFile struct {
	Status string
	File   string
}

type HistoryFiles []*HistoryFile
