package result

import (
	"time"
)

type Result struct {
	ID  string    `json:"id"`
	Src string    `json:"src"`
	Tgt string    `json:"tgt"`
	Act []string  `json:"act"`
	At  time.Time `json:"at"`
}

func (r *Result) String() string {
	return r.ID
}
