package cutil_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
)

func TestPathHelpers(t *testing.T) {
	u := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/items/123", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id":   "123",
		"flag": "true",
		"arr":  "a, b",
		"uuid": u.String(),
	})

	if got, err := cutil.PathString(req, "id", false); err != nil || got != "123" {
		t.Fatalf("PathString returned [%v], err=[%v]", got, err)
	}

	if _, err := cutil.PathString(req, "missing", false); err == nil {
		t.Fatalf("expected error for missing path var")
	}

	if got, err := cutil.PathBool(req, "flag"); err != nil || got != true {
		t.Fatalf("PathBool returned [%v], err=[%v]", got, err)
	}

	if got, err := cutil.PathInt(req, "id"); err != nil || got != 123 {
		t.Fatalf("PathInt returned [%v], err=[%v]", got, err)
	}

	if got, err := cutil.PathArray(req, "arr"); err != nil || !reflect.DeepEqual(got, util.Strings{"a", "b"}) {
		t.Fatalf("PathArray returned [%v], err=[%v]", got, err)
	}

	gotUUID, err := cutil.PathUUID(req, "uuid")
	if err != nil || gotUUID == nil || *gotUUID != u {
		t.Fatalf("PathUUID returned [%v], err=[%v]", gotUUID, err)
	}
}

func TestQueryStringHelpers(t *testing.T) {
	u := uuid.New()
	uri, err := url.Parse("http://example.com/?a=1&b=true&c=foo&uuid=" + u.String() + "&multi=a&multi=b")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if got := cutil.QueryStringString(uri, "c"); got != "foo" {
		t.Fatalf("QueryStringString returned [%v]", got)
	}

	if got := cutil.QueryStringBool(uri, "b"); !got {
		t.Fatalf("QueryStringBool returned false")
	}

	if got := cutil.QueryStringInt(uri, "a"); got != 1 {
		t.Fatalf("QueryStringInt returned [%v]", got)
	}

	if got := cutil.QueryStringUUID(uri, "uuid"); got == nil || *got != u {
		t.Fatalf("QueryStringUUID returned [%v]", got)
	}

	m := cutil.QueryStringAsMap(uri)
	if !reflect.DeepEqual(m["multi"], []string{"a", "b"}) {
		t.Fatalf("QueryStringAsMap returned [%v]", m["multi"])
	}
}

func TestHeaderMaps(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	req.Header.Add("X-Test", "a")
	req.Header.Add("X-Multi", "a")
	req.Header.Add("X-Multi", "b")

	rm := cutil.RequestHeadersMap(req)
	if rm["X-Test"] != "a" {
		t.Fatalf("RequestHeadersMap returned [%v]", rm["X-Test"])
	}
	if !reflect.DeepEqual(rm["X-Multi"], []string{"a", "b"}) {
		t.Fatalf("RequestHeadersMap returned [%v]", rm["X-Multi"])
	}

	rr := httptest.NewRecorder()
	rr.Header().Add("X-Test", "a")
	rr.Header().Add("X-Multi", "a")
	rr.Header().Add("X-Multi", "b")

	wm := cutil.ResponseHeadersMap(rr)
	if wm["X-Test"] != "a" {
		t.Fatalf("ResponseHeadersMap returned [%v]", wm["X-Test"])
	}
	if !reflect.DeepEqual(wm["X-Multi"], []string{"a", "b"}) {
		t.Fatalf("ResponseHeadersMap returned [%v]", wm["X-Multi"])
	}
}
