package cutil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"projectforge.dev/projectforge/app/controller/cutil"
)

func TestWriteCounterTracksBytesAndStatus(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	wc := cutil.NewWriteCounter(rr)

	if wc.Started().IsZero() {
		t.Fatalf("expected started time to be set")
	}

	wc.WriteHeader(http.StatusTeapot)
	if wc.StatusCode() != http.StatusTeapot {
		t.Fatalf("StatusCode was [%d]", wc.StatusCode())
	}
	if rr.Header().Get("X-Runtime") == "" {
		t.Fatalf("X-Runtime header missing")
	}

	if _, err := wc.Write([]byte("abc")); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if wc.Count() != 3 {
		t.Fatalf("Count was [%d]", wc.Count())
	}

	if wc.Unwrap() != rr {
		t.Fatalf("Unwrap did not return original ResponseWriter")
	}
}

func TestWriteCounterHijackUnsupported(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	wc := cutil.NewWriteCounter(rr)
	if _, _, err := wc.Hijack(); err == nil {
		t.Fatalf("expected hijack error for non-hijacker response")
	}
}
