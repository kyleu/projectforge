package cutil_test

import (
	"net/url"
	"reflect"
	"testing"

	"projectforge.dev/projectforge/app/controller/cutil"
)

func TestAddRoute(t *testing.T) {
	t.Parallel()
	orig := cutil.AppRoutesList
	cutil.AppRoutesList = map[string][]string{}
	defer func() { cutil.AppRoutesList = orig }()

	cutil.AddRoute("GET", "/b")
	cutil.AddRoute("GET", "/a")
	cutil.AddRoute("GET", "/a")

	if got := cutil.AppRoutesList["GET"]; !reflect.DeepEqual(got, []string{"/a", "/b"}) {
		t.Fatalf("AddRoute returned [%v]", got)
	}
}

func TestURLAddQuery(t *testing.T) {
	t.Parallel()
	u, _ := url.Parse("http://example.com/?a=1")
	cutil.URLAddQuery(u, "b", "2")

	if got := u.Query().Get("a"); got != "1" {
		t.Fatalf("query a was [%v]", got)
	}
	if got := u.Query().Get("b"); got != "2" {
		t.Fatalf("query b was [%v]", got)
	}
}
