// Content managed by Project Forge, see [projectforge.md] for details.
package auth

import (
	"github.com/valyala/fasthttp"
)

type params struct {
	q *fasthttp.Args
}

func (p *params) Get(key string) string {
	b := p.q.Peek(key)
	if len(b) > 0 {
		return string(b)
	}
	return ""
}
