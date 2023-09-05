package har

import (
	"net/url"

	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type PostParam struct {
	Name        string `json:"name"`
	Value       string `json:"value,omitempty"`
	FileName    string `json:"fileName,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

type PostParams []*PostParam

type PostData struct {
	MimeType string     `json:"mimeType"`
	Params   PostParams `json:"params"`
	Text     string     `json:"text"`
	Comment  string     `json:"comment,omitempty"`
}

func (p *PostData) String() string {
	var ret string
	if len(p.Params) > 0 {
		ret += util.StringPlural(len(p.Params), "param")
	}
	if p.MimeType != "" {
		ret += " (" + p.MimeType + ")"
	}
	return ret
}

type Request struct {
	Method      string    `json:"method"`
	URL         string    `json:"url"`
	HTTPVersion string    `json:"httpVersion"`
	Cookies     Cookies   `json:"cookies"`
	Headers     NVPs      `json:"headers"`
	QueryString NVPs      `json:"queryString"`
	PostData    *PostData `json:"postData"`
	HeadersSize int       `json:"headersSize"`
	BodySize    int       `json:"bodySize"`
	Comment     string    `json:"comment"`
}

func (r *Request) Size() int {
	if r.BodySize > 0 {
		return r.HeadersSize + r.BodySize
	}
	if r.PostData == nil {
		return r.HeadersSize
	}
	return r.HeadersSize + len(r.PostData.Text)
}

func (r *Request) GetURL() *url.URL {
	u, err := url.Parse(r.URL)
	if err != nil {
		return nil
	}
	return u
}

func (r *Request) WithReplacements(repl func(s string) string) *Request {
	cooks := lo.Map(r.Cookies, func(c *Cookie, _ int) *Cookie {
		return &Cookie{
			Name: repl(c.Name), Value: repl(c.Value), Path: repl(c.Path), Domain: repl(c.Domain),
			Expires: c.Expires, HTTPOnly: c.HTTPOnly, Secure: c.Secure, Comment: c.Comment,
		}
	})
	var pd *PostData
	if r.PostData != nil {
		var pp PostParams
		for _, p := range r.PostData.Params {
			pp = append(pp, &PostParam{Name: repl(p.Name), Value: repl(p.Value), FileName: repl(p.FileName), ContentType: p.ContentType, Comment: p.Comment})
		}
		pd = &PostData{MimeType: repl(r.PostData.MimeType), Params: pp, Text: repl(r.PostData.Text), Comment: repl(r.PostData.Comment)}
	}
	return &Request{
		Method: r.Method, URL: repl(r.URL), HTTPVersion: r.HTTPVersion, Cookies: cooks, Headers: r.Headers.WithReplacements(repl),
		QueryString: r.QueryString.WithReplacements(repl), PostData: pd, HeadersSize: r.HeadersSize, BodySize: r.BodySize, Comment: repl(r.Comment),
	}
}
