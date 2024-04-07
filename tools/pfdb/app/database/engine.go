package database

import "fmt"

type DBType struct {
	Key               string `json:"key"`
	Title             string `json:"title"`
	Quote             string `json:"-"`
	Placeholder       string `json:"-"`
	SupportsReturning bool   `json:"-"`
}

func (t *DBType) PlaceholderFor(idx int) string {
	switch t.Placeholder {
	case "$", "":
		return fmt.Sprintf("$%d", idx)
	case "@":
		return fmt.Sprintf("@p%d", idx)
	default:
		return t.Placeholder
	}
}

func (t *DBType) Quoted(s string) string {
	return fmt.Sprintf("%s%s%s", t.Quote, s, t.Quote)
}
