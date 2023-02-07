// Content managed by Project Forge, see [projectforge.md] for details.
package cutil

import (
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/lib/filter"
)

func ParamSetFromRequest(rc *fasthttp.RequestCtx) filter.ParamSet {
	ret := filter.ParamSet{}
	args := rc.URI().QueryArgs()
	args.VisitAll(func(key []byte, value []byte) {
		qk := string(key)
		if strings.Contains(qk, ".") {
			ret = apply(ret, qk, string(args.Peek(qk)))
		}
	})
	return ret
}

func apply(ps filter.ParamSet, qk string, qv string) filter.ParamSet {
	switch {
	case strings.HasSuffix(qk, ".o"):
		curr := getCurr(ps, strings.TrimSuffix(qk, ".o"))
		asc := true
		if strings.HasSuffix(qv, ".d") {
			asc = false
			qv = qv[0 : len(qv)-2]
		}
		curr.Orderings = append(curr.Orderings, &filter.Ordering{Column: qv, Asc: asc})
	case strings.HasSuffix(qk, ".l"):
		curr := getCurr(ps, strings.TrimSuffix(qk, ".l"))
		li, err := strconv.ParseInt(qv, 10, 64)
		if err == nil {
			curr.Limit = int(li)
			max := 1000
			if curr.Limit > max {
				curr.Limit = max
			}
		}
	case strings.HasSuffix(qk, ".x"):
		curr := getCurr(ps, strings.TrimSuffix(qk, ".x"))
		xi, err := strconv.ParseInt(qv, 10, 64)
		if err == nil {
			curr.Offset = int(xi)
		}
	}
	return ps
}

func getCurr(q filter.ParamSet, key string) *filter.Params {
	curr, ok := q[key]
	if !ok {
		curr = &filter.Params{Key: key}
		q[key] = curr
	}
	return curr
}
