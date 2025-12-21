package queue

import "{{{ .Package }}}/app/util"

type Status struct {
	Name     string         `json:"name"`
	Limit    int            `json:"limit,omitzero"`
	Timeout  string         `json:"timeout,omitzero"`
	Table    string         `json:"table,omitzero"`
	Started  string         `json:"started"`
	Sent     map[string]int `json:"sent,omitzero"`
	Received map[string]int `json:"received,omitzero"`
}

func (q *Queue) Status() *Status {
	return &Status{
		Name:     q.name,
		Limit:    q.limit,
		Timeout:  util.MicrosToMillis(int(q.timeout / 1000)),
		Table:    q.table,
		Started:  util.TimeRelative(&q.started),
		Sent:     util.MapClone(q.sent),
		Received: util.MapClone(q.received),
	}
}
