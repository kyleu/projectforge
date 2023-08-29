package har

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/net/http/httpguts"

	"{{{ .Package }}}/app/util"
)

type Entry struct {
	PageRef         string       `json:"pageref,omitempty"`
	StartedDateTime string       `json:"startedDateTime"`
	Time            float32      `json:"time"`
	Request         *Request     `json:"request"`
	Response        *Response    `json:"response"`
	Cache           *Cache       `json:"cache"`
	PageTimings     *PageTimings `json:"timings"`
	ServerIPAddress string       `json:"serverIPAddress,omitempty"`
	Connection      string       `json:"connection,omitempty"`
	Comment         string       `json:"comment,omitempty"`
	Selector        *Selector    `json:"-"`
}

func (e *Entry) String() string {
	if len(e.Request.URL) > 64 {
		return e.Request.Method + " " + e.Request.URL[:64] + "..."
	}
	return e.Request.Method + " " + e.Request.URL
}

func (e *Entry) Duration() int {
	return e.PageTimings.Wait * 1000
}

func (e *Entry) Clone() *Entry {
	return &Entry{
		PageRef: e.PageRef, StartedDateTime: e.StartedDateTime, Time: e.Time, Request: e.Request, Response: e.Response, Cache: e.Cache,
		PageTimings: e.PageTimings, ServerIPAddress: e.ServerIPAddress, Connection: e.Connection, Comment: e.Comment, Selector: e.Selector,
	}
}

func (e *Entry) Cleaned() *Entry {
	ret := e
	if ret.Request.PostData != nil && len(ret.Request.PostData.Text) > 1024*16 {
		ret = ret.Clone()
		ret.Request.PostData.Text = util.ByteSizeSI(int64(len(ret.Request.PostData.Text)))
	}
	if ret.Response != nil && ret.Response.Content != nil && ret.Response.Content.Size > 1024*16 {
		ret = ret.Clone()
		ret.Response.Content.Text = util.ByteSizeSI(int64(ret.Response.Content.Size))
	}
	return ret
}

func (e *Entry) ToRequest(ctx context.Context, ignoreCookies bool) (*http.Request, error) {
	body := ""

	if e.Request.PostData != nil {
		if len(e.Request.PostData.Params) == 0 {
			body = e.Request.PostData.Text
		} else {
			form := url.Values{}
			for _, p := range e.Request.PostData.Params {
				form.Add(p.Name, p.Value)
			}
			body = form.Encode()
		}
	}

	req, err := http.NewRequestWithContext(ctx, e.Request.Method, e.Request.URL, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return nil, err
	}

	for _, h := range e.Request.Headers {
		if httpguts.ValidHeaderFieldName(h.Name) && httpguts.ValidHeaderFieldValue(h.Value) && h.Name != "Cookie" {
			req.Header.Add(h.Name, h.Value)
		}
	}

	if !ignoreCookies {
		for _, c := range e.Request.Cookies {
			cookie := &http.Cookie{Name: c.Name, Value: url.QueryEscape(c.Value), HttpOnly: false, Domain: c.Domain}
			req.AddCookie(cookie)
		}
	}

	return req, nil
}

func (e *Entry) WithReplacementsMap(repls map[string]string, vars util.ValueMap) *Entry {
	if len(repls) == 0 {
		return e
	}
	startToken, endToken := "{{", "}}"
	return e.WithReplacements(func(s string) string {
		for k, v := range repls {
			if v == "" {
				v = k
			}
			if strings.Contains(v, "||") {
				v = util.StringSplitAndTrim(v, "||")[0]
			}
			sIdx := strings.Index(v, startToken)
			for sIdx > -1 {
				eIdx := strings.Index(v, endToken)
				if eIdx == -1 {
					return "missing end token [" + endToken + "]"
				}
				match := v[sIdx+3 : eIdx]
				r := startToken + match + endToken
				variable := vars.GetStringOpt(strings.TrimSpace(match))
				if variable == "" {
					return "missing variable [" + strings.TrimSpace(match) + "]"
				}
				v = strings.Replace(v, r, variable, 1)
				sIdx = strings.Index(v, startToken)
			}
			s = strings.ReplaceAll(s, k, v)
		}
		return s
	})
}

func (e *Entry) WithReplacements(repls func(s string) string) *Entry {
	ret := e.Clone()
	ret.Request = ret.Request.WithReplacements(repls)
	ret.Response = ret.Response.WithReplacements(repls)
	return ret
}

type Entries []*Entry

func (e Entries) ForPage(ref string) Entries {
	return lo.Filter(e, func(x *Entry, _ int) bool {
		return x.PageRef == ref
	})
}

func (e Entries) BySelector(sel *Selector) Entries {
	return lo.Filter(e, func(x *Entry, _ int) bool {
		return x.Selector != nil && x.Selector.Matches(sel)
	})
}

func (e Entries) Selectors() Selectors {
	return lo.UniqBy(lo.Map(e, func(x *Entry, _ int) *Selector {
		return x.Selector
	}), func(s *Selector) string {
		return s.String()
	})
}

func (e Entries) Trimmed() Entries {
	return lo.Filter(e, func(x *Entry, _ int) bool {
		return x.Response != nil && x.Response.Status != 0
	})
}

func (e Entries) TotalDuration() int {
	return lo.SumBy(e, func(x *Entry) int {
		return x.Duration()
	})
}

func (e Entries) TotalResponseBodySize() int {
	return lo.SumBy(e, func(x *Entry) int {
		if x.Response == nil || x.Response.Content == nil {
			return 0
		}
		return x.Response.Content.Size
	})
}

func (e Entries) WithReplacementsMap(repls map[string]string, vars util.ValueMap) Entries {
	if len(repls) == 0 {
		return e
	}
	return e.WithReplacementsMap(repls, vars)
}

func (e Entries) WithReplacements(repl func(s string) string) Entries {
	return lo.Map(e, func(x *Entry, _ int) *Entry {
		return x.WithReplacements(repl)
	})
}

func (e Entries) Cleaned() Entries {
	return lo.Map(e, func(x *Entry, _ int) *Entry {
		return x.Cleaned()
	})
}
