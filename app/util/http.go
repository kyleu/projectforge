// Package util - Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"context"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type HTTPRequest struct {
	*http.Request
	client *http.Client
}

func NewHTTPRequest(ctx context.Context, method string, url string) *HTTPRequest {
	if method == "" {
		method = http.MethodGet
	}
	ret, _ := http.NewRequestWithContext(ctx, method, url, http.NoBody)
	return &HTTPRequest{Request: ret}
}

func (r *HTTPRequest) WithHeader(k string, v string) *HTTPRequest {
	r.Header.Set(k, v)
	return r
}

func (r *HTTPRequest) WithClient(c *http.Client) *HTTPRequest {
	r.client = c
	return r
}

func (r *HTTPRequest) RunSimple() (*http.Response, []byte, error) {
	cl := Choose(r.client == nil, http.DefaultClient, r.client)
	rsp, err := cl.Do(r.Request)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = rsp.Body.Close() }()
	b, _ := io.ReadAll(rsp.Body)
	if rsp.StatusCode != http.StatusOK {
		return rsp, b, errors.Errorf("response from url [%s] has status [%d], expected [200]", r.URL, rsp.StatusCode)
	}
	return rsp, b, nil
}
