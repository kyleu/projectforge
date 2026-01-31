//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"{{{ .Package }}}/app/util"
)

func TestHTTPRequest(t *testing.T) {
	var gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Test")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	req := util.NewHTTPRequest(context.Background(), http.MethodGet, srv.URL)
	req.WithHeader("X-Test", "1").WithClient(srv.Client())

	status, _, body, err := req.RunSimple()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != http.StatusOK || string(body) != "ok" {
		t.Fatalf("unexpected response: %d %s", status, string(body))
	}
	if gotHeader != "1" {
		t.Fatalf("expected header to be set, got %q", gotHeader)
	}

	bad := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer bad.Close()

	_, _, _, err = util.NewHTTPGet(context.Background(), bad.URL).WithClient(bad.Client()).RunSimple()
	if err == nil || !strings.Contains(err.Error(), "expected [200]") {
		t.Fatalf("expected RunSimple error for non-200")
	}
}
