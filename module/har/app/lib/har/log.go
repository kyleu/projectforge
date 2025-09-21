package har

import "net/url"

type Creator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Comment string `json:"comment,omitzero"`
}

type Browser struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Comment string `json:"comment"`
}

type Log struct {
	Key     string   `json:"-"`
	Version string   `json:"version"`
	Creator *Creator `json:"creator"`
	Browser *Browser `json:"browser"`
	Pages   Pages    `json:"pages,omitempty"`
	Entries Entries  `json:"entries"`
	Comment string   `json:"comment"`
}

func (l *Log) WebPath() string {
	return "/har/" + url.QueryEscape(l.Key)
}

type Logs []*Log
