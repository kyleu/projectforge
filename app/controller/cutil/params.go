package cutil

import (
	"net/http"
	"strconv"
	"strings"

	"projectforge.dev/projectforge/app/lib/filter"
)

func ParamSetFromRequest(r *http.Request) filter.ParamSet {
	ret := filter.ParamSet{}
	for k, v := range r.URL.Query() {
		if strings.Contains(k, ".") {
			ret = apply(ret, k, strings.Join(v, ","))
		}
	}
	return ret
}

func apply(ps filter.ParamSet, qk string, qv string) filter.ParamSet {
	switch {
	case strings.HasSuffix(qk, filter.SuffixOrder):
		curr := getCurr(ps, strings.TrimSuffix(qk, filter.SuffixOrder))
		asc := true
		if strings.HasSuffix(qv, filter.SuffixDescending) {
			asc = false
			qv = qv[0 : len(qv)-2]
		}
		curr.Orderings = append(curr.Orderings, &filter.Ordering{Column: qv, Asc: asc})
	case strings.HasSuffix(qk, filter.SuffixLimit):
		curr := getCurr(ps, strings.TrimSuffix(qk, filter.SuffixLimit))
		li, err := strconv.ParseInt(qv, 10, 32)
		if err == nil {
			curr.Limit = int(li)
			maxCount := 100000
			if curr.Limit > maxCount {
				curr.Limit = maxCount
			}
		}
	case strings.HasSuffix(qk, filter.SuffixOffset):
		curr := getCurr(ps, strings.TrimSuffix(qk, filter.SuffixOffset))
		xi, err := strconv.ParseInt(qv, 10, 32)
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
