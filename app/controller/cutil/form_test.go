package cutil_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"projectforge.dev/projectforge/app/controller/cutil"
)

func TestParseFormJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(`{"a":"b"}`))
	req.Header.Set(cutil.HeaderContentType, "application/json")

	m, err := cutil.ParseForm(req, []byte(`{"a":"b"}`))
	if err != nil {
		t.Fatalf("ParseForm returned error: %v", err)
	}
	if m.GetStringOpt("a") != "b" {
		t.Fatalf("ParseForm returned [%v]", m)
	}
}

func TestParseFormAsMapsJSONArray(t *testing.T) {
	body := `[{"a":1},{"a":2}]`
	req := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(body))
	req.Header.Set(cutil.HeaderContentType, "application/json")

	m, err := cutil.ParseFormAsMaps(req, []byte(body))
	if err != nil {
		t.Fatalf("ParseFormAsMaps returned error: %v", err)
	}
	if len(m) != 2 {
		t.Fatalf("ParseFormAsMaps returned %d maps", len(m))
	}
}

func TestParseFormHTTPFallback(t *testing.T) {
	body := "a=1&b=two"
	req := httptest.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(body))
	req.Header.Set(cutil.HeaderContentType, "application/x-www-form-urlencoded")

	m, err := cutil.ParseForm(req, []byte(body))
	if err != nil {
		t.Fatalf("ParseForm returned error: %v", err)
	}
	if m.GetStringOpt("a") != "1" || m.GetStringOpt("b") != "two" {
		t.Fatalf("ParseForm returned [%v]", m)
	}
}

func TestCleanID(t *testing.T) {
	if got := cutil.CleanID("key", "id"); got != "id" {
		t.Fatalf("CleanID returned [%v]", got)
	}
	if got := cutil.CleanID("key", ""); !strings.HasPrefix(got, "key-") {
		t.Fatalf("CleanID returned [%v]", got)
	}
}
