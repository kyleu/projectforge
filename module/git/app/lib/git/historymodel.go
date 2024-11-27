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

type HistoryAuthor struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Count int    `json:"count"`
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

func (h HistoryEntries) Authors() []*HistoryAuthor {
	byAuthor := lo.GroupBy(h, func(x *HistoryEntry) string {
		return x.AuthorEmail
	})
	ret := make([]*HistoryAuthor, 0, len(byAuthor))
	for k, v := range byAuthor {
		ret = append(ret, &HistoryAuthor{
			Key:   k,
			Name:  v[0].AuthorName,
			Count: len(v),
		})
	}
	return ret
}

type HistoryFile struct {
	Status string
	File   string
}

func (h *HistoryFile) String() string {
	return h.File
}

type HistoryFiles []*HistoryFile

func (h HistoryFiles) Strings() []string {
	return lo.Map(h, func(x *HistoryFile, _ int) string {
		return x.String()
	})
}
