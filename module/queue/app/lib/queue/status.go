package queue

import (
	"maps"

	"{{{ .Package }}}/app/util"
)

type Status struct {
	Name     string         `json:"name"`
	Limit    int            `json:"limit,omitempty"`
	Timeout  string         `json:"timeout,omitempty"`
	Table    string         `json:"table,omitempty"`
	Started  string         `json:"started"`
	Sent     map[string]int `json:"sent,omitempty"`
	Received map[string]int `json:"received,omitempty"`
}

func (q *Queue) Status() *Status {
	return &Status{
		Name:     q.name,
		Limit:    q.limit,
		Timeout:  util.MicrosToMillis(int(q.timeout / 1000)),
		Table:    q.table,
		Started:  util.TimeRelative(&q.started),
		Sent:     maps.Clone(q.sent),
		Received: maps.Clone(q.received),
	}
}
