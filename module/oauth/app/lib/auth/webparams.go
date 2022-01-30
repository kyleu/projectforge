package auth

import (
	"github.com/valyala/fasthttp"
)

type params struct {
	q *fasthttp.Args
}

func (p *params) Get(key string) string {
	if b := p.q.Peek(key); len(b) > 0 {
		return string(b)
	}
	return ""
}
