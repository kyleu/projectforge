package har

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Selector struct {
	Har     string `json:"har,omitempty"`
	URL     string `json:"url,omitempty"`
	Mime    string `json:"mime,omitempty"`
	Idx     int    `json:"idx,omitempty"`
	Comment string `json:"comment,omitempty"`
}

func (s Selector) Matches(x *Selector) bool {
	return s.Har == x.Har && s.URL == x.URL && s.Mime == x.Mime && s.Idx == x.Idx
}

func (s Selector) String() string {
	var ret []string
	if s.Har != "" {
		ret = append(ret, "Archive: "+s.Har)
	}
	if s.URL != "" {
		ret = append(ret, "URL: "+s.URL)
	}
	if s.Mime != "" {
		ret = append(ret, "MIME: "+s.Mime)
	}
	if s.Idx != 0 {
		ret = append(ret, fmt.Sprintf("Index: %d", s.Idx))
	}
	if len(ret) == 0 {
		ret = append(ret, "*")
	}
	return util.StringJoin(ret, ", ")
}

type Selectors []*Selector

func (e Entries) Find(s *Selector) (Entries, error) {
	matches := func(k string, v string) bool {
		pre, suff := strings.HasPrefix(k, "*"), strings.HasSuffix(k, "*")
		k = strings.TrimSuffix(strings.TrimPrefix(k, "*"), "*")
		if pre && suff {
			return strings.Contains(v, k)
		}
		if pre {
			return strings.HasSuffix(v, k)
		}
		if suff {
			return strings.HasPrefix(v, k)
		}
		return k == v
	}
	ret := util.ArrayCopy(e)
	if s.URL != "" && s.URL != "*" {
		ret = lo.Filter(ret, func(e *Entry, _ int) bool {
			return matches(s.URL, e.Request.URL)
		})
	}
	if s.Mime != "" && s.Mime != "*" {
		ret = lo.Filter(ret, func(e *Entry, _ int) bool {
			if e.Response.Content == nil {
				return false
			}
			tst := e.Response.Content.MimeType
			if idx := strings.Index(tst, ";"); idx > -1 {
				tst = tst[:idx]
			}
			return matches(s.Mime, tst)
		})
	}
	if s.Idx > 0 {
		if s.Idx > len(ret) {
			return nil, errors.Errorf("index [%d] does not exist among [%d] entries", s.Idx-1, len(ret))
		}
		return Entries{ret[s.Idx-1]}, nil
	}
	return ret, nil
}
