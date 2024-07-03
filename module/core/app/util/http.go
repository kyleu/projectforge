package util

import (
	"context"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

var HTTPDefaultClient = http.DefaultClient

type HTTPRequest struct {
	*http.Request
	client *http.Client
}

func NewHTTPRequest(ctx context.Context, method string, url string, bodies ...io.Reader) *HTTPRequest {
	if method == "" {
		method = http.MethodGet
	}
	var body io.Reader = http.NoBody
	if len(bodies) > 0 {
		body = bodies[0]
	}
	ret, _ := http.NewRequestWithContext(ctx, method, url, body)
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

func (r *HTTPRequest) Run() (*http.Response, error) {
	cl := Choose(r.client == nil, HTTPDefaultClient, r.client)
	rsp, err := cl.Do(r.Request)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (r *HTTPRequest) RunSimple() (*http.Response, []byte, error) {
	rsp, err := r.Run()
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
