package cutil_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"{{{ .Package }}}/app/controller/cutil"
)

func TestWriteCORSUsesReferer(t *testing.T) {
	rr := httptest.NewRecorder()
	rr.Header().Set(cutil.HeaderReferer, "http://example.com:8080/path")

	cutil.WriteCORS(rr)

	if got := rr.Header().Get(cutil.HeaderAccessControlAllowOrigin); got != "http://example.com:8080" {
		t.Fatalf("Access-Control-Allow-Origin was [%v]", got)
	}
	if got := rr.Header().Get(cutil.HeaderAccessControlAllowMethods); got == "" {
		t.Fatalf("Access-Control-Allow-Methods missing")
	}
	if got := rr.Header().Get(cutil.HeaderAccessControlAllowHeaders); got == "" {
		t.Fatalf("Access-Control-Allow-Headers missing")
	}
}

func TestGetContentTypesFormatOverride(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com/?format=csv", nil)
	req.Header.Set(cutil.HeaderAccept, "application/json")

	ct, format := cutil.GetContentTypes(req)
	if !cutil.IsContentTypeCSV(ct) || format != "csv" {
		t.Fatalf("GetContentTypes returned [%v], [%v]", ct, format)
	}

	req2 := httptest.NewRequest(http.MethodGet, "http://example.com/", nil)
	req2.Header.Set(cutil.HeaderAccept, "text/html; charset=UTF-8")
	ct2, format2 := cutil.GetContentTypes(req2)
	if ct2 != "text/html" || format2 != "" {
		t.Fatalf("GetContentTypes returned [%v], [%v]", ct2, format2)
	}
}

func TestRespondMIMEWritesHeadersAndBody(t *testing.T) {
	rr := httptest.NewRecorder()
	wc := cutil.NewWriteCounter(rr)
	payload := []byte("hello")

	if _, err := cutil.RespondMIME("file.txt", "text/plain", payload, wc); err != nil {
		t.Fatalf("RespondMIME returned error: %v", err)
	}

	if got := rr.Header().Get(cutil.HeaderContentType); !strings.HasPrefix(got, "text/plain") || !strings.Contains(got, "charset=") {
		t.Fatalf("Content-Type was [%v]", got)
	}

	if got := rr.Header().Get("Content-Disposition"); !strings.Contains(got, "file.txt") {
		t.Fatalf("Content-Disposition was [%v]", got)
	}

	if got := rr.Header().Get(cutil.HeaderAccessControlAllowOrigin); got != "*" {
		t.Fatalf("Access-Control-Allow-Origin was [%v]", got)
	}

	if body := rr.Body.String(); body != "hello" {
		t.Fatalf("body was [%v]", body)
	}

	if wc.Count() != int64(len(payload)) {
		t.Fatalf("WriteCounter count was [%d]", wc.Count())
	}
}

func TestRespondMIMERejectsEmptyPayload(t *testing.T) {
	rr := httptest.NewRecorder()
	wc := cutil.NewWriteCounter(rr)
	if _, err := cutil.RespondMIME("", "text/plain", nil, wc); err == nil {
		t.Fatalf("expected error for empty payload")
	}
}

func TestWriteCORSNetworkScheme(t *testing.T) {
	rr := httptest.NewRecorder()
	rr.Header().Set(cutil.HeaderReferer, "http://example.network/path")
	cutil.WriteCORS(rr)

	if got := rr.Header().Get(cutil.HeaderAccessControlAllowOrigin); got != "https://example.network" {
		t.Fatalf("Access-Control-Allow-Origin was [%v]", got)
	}
}

func TestGetContentTypesFromContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://example.com/?t=yml", nil)
	req.Header.Set(cutil.HeaderContentType, "application/x-yaml")

	ct, format := cutil.GetContentTypes(req)
	if !cutil.IsContentTypeYAML(ct) || format != "yml" {
		t.Fatalf("GetContentTypes returned [%v], [%v]", ct, format)
	}
}

func TestGetContentTypesFromQueryParam(t *testing.T) {
	u, _ := url.Parse("http://example.com/?t=debug")
	req := &http.Request{URL: u, Header: http.Header{}}
	ct, format := cutil.GetContentTypes(req)
	if !cutil.IsContentTypeDebug(ct) || format != "debug" {
		t.Fatalf("GetContentTypes returned [%v], [%v]", ct, format)
	}
}
