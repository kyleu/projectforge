package git

import (
	"time"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type HistoryArgs struct {
	Path    string     `json:"path,omitempty"`
	Since   *time.Time `json:"since,omitempty"`
	Authors []string   `json:"authors,omitempty"`
	Limit   int        `json:"limit,omitempty"`
	Commit  string     `json:"commit,omitempty"`
	Debug   bool       `json:"debug,omitempty"`
}

type HistoryResult struct {
	Args    *HistoryArgs   `json:"args,omitempty"`
	Entries HistoryEntries `json:"entries"`
	Debug   any            `json:"debug,omitempty"`
}

type HistoryEntry struct {
	Headers     util.ValueMap `json:"headers" xml:"headers"`
	SHA         string        `json:"sha" xml:"sha"`
	AuthorName  string        `json:"authorName" xml:"authorName"`
	AuthorEmail string        `json:"authorEmail" xml:"authorEmail"`
	Message     string        `json:"message" xml:"message"`
	Occurred    time.Time     `json:"occurred" xml:"occurred"`
	Files       HistoryFiles  `json:"files" xml:"files"`
}

type HistoryEntries []*HistoryEntry

func (h HistoryEntries) Get(sha string) *HistoryEntry {
	return lo.FindOrElse(h, nil, func(x *HistoryEntry) bool {
		return x.SHA == sha
	})
}

type HistoryFile struct {
	Status string
	File   string
}

type HistoryFiles []*HistoryFile
