// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import "fmt"

type Progress struct {
	Key       string `json:"key,omitempty"`
	Total     int    `json:"total,omitempty"`
	Completed int    `json:"completed,omitempty"`
	fns       []func(p *Progress, delta int)
}

func NewProgress(key string, total int, fns ...func(p *Progress, delta int)) *Progress {
	if total == 0 {
		total = 100
	}
	return &Progress{Key: key, Total: total, fns: fns}
}

func (p *Progress) String() string {
	return fmt.Sprintf("%s (%d of %d)", p.Key, p.Completed, p.Total)
}

func (p *Progress) Increment(i int, logger Logger) {
	if p == nil {
		return
	}
	p.Completed += i
	p.call(i, logger)
}

func (p *Progress) call(delta int, logger Logger) {
	if logger != nil {
		logger.Debugf("%s [%d]", p.String(), delta)
	}
	for _, fn := range p.fns {
		fn(p, delta)
	}
}
