package har

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Response struct {
	Status      int      `json:"status"`
	StatusText  string   `json:"statusText"`
	HTTPVersion string   `json:"httpVersion"`
	Cookies     Cookies  `json:"cookies"`
	Headers     NVPs     `json:"headers"`
	Content     *Content `json:"content"`
	RedirectURL string   `json:"redirectURL"`
	HeadersSize int      `json:"headersSize"`
	BodySize    int      `json:"bodySize"`
	Comment     string   `json:"comment,omitempty"`
}

func ResponseFromHTTP(r *http.Response) (*Response, error) {
	cooks := lo.Map(r.Cookies(), func(c *http.Cookie, _ int) *Cookie {
		exp := c.Expires.Format("2006-01-02T15:04:05.000Z")
		return &Cookie{Name: c.Name, Value: c.Value, Path: c.Path, Domain: c.Domain, Expires: exp, HTTPOnly: c.HttpOnly, Secure: c.Secure}
	})
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var headers NVPs
	for k, vs := range r.Header {
		for _, v := range vs {
			headers = append(headers, &NVP{Name: k, Value: v})
		}
	}
	enc := headers.GetValue("Content-Encoding")
	switch enc {
	case "gzip":
		zr, err := gzip.NewReader(bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		defer func() { _ = zr.Close() }()
		bodyBytes, err = io.ReadAll(zr)
		if err != nil {
			return nil, err
		}
	case "br":
		br := brotli.NewReader(bytes.NewReader(bodyBytes))
		bodyBytes, err = io.ReadAll(br)
		if err != nil {
			return nil, err
		}
	}
	var js any
	_ = util.FromJSON(bodyBytes, &js)
	content := &Content{Size: len(bodyBytes), Text: string(bodyBytes), JSON: js}
	ret := &Response{Status: r.StatusCode, StatusText: r.Status, Cookies: cooks, Headers: headers, Content: content, BodySize: content.Size}
	return ret, nil
}

func (r *Response) Size() int {
	if r.Content == nil {
		return 0
	}
	if r.Content.Size != 0 {
		return r.Content.Size
	}
	return len(r.Content.Text)
}

func (r *Response) BodyString() string {
	if r.Content == nil {
		return ""
	}
	return r.Content.Text
}

func (r *Response) ContentType() string {
	ret := r.Headers.GetValue("content-type")
	if idx := strings.LastIndex(ret, ";"); idx > -1 {
		ret = ret[:idx]
	}
	return ret
}

func (r *Response) WithReplacements(repl func(s string) string) *Response {
	c := &Content{
		Size: r.Content.Size, Compression: r.Content.Compression, MimeType: repl(r.Content.MimeType),
		Text: repl(r.Content.Text), Encoding: r.Content.Encoding, Comment: r.Content.Comment, File: repl(r.Content.File),
	}
	_ = util.FromJSON([]byte(r.Content.Text), &c.JSON)
	return &Response{
		Status:      r.Status,
		StatusText:  r.StatusText,
		HTTPVersion: r.HTTPVersion,
		Cookies:     r.Cookies,
		Headers:     r.Headers.WithReplacements(repl),
		Content:     c,
		RedirectURL: repl(r.RedirectURL),
		HeadersSize: r.HeadersSize,
		BodySize:    r.BodySize,
		Comment:     r.Comment,
	}
}
