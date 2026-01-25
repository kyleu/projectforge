package cutil_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/filter"
)

func TestParamSetFromRequest(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/?user.o=name&user.l=200001&user.x=5&order.o=title.d", http.NoBody)

	ps := cutil.ParamSetFromRequest(req)
	user := ps["user"]
	if user == nil {
		t.Fatalf("missing user params")
	}

	if user.Limit != 100000 {
		t.Fatalf("expected limit 100000, got %d", user.Limit)
	}
	if user.Offset != 5 {
		t.Fatalf("expected offset 5, got %d", user.Offset)
	}
	if len(user.Orderings) != 1 || user.Orderings[0].Column != "name" || !user.Orderings[0].Asc {
		t.Fatalf("unexpected user orderings: %+v", user.Orderings)
	}

	order := ps["order"]
	if order == nil {
		t.Fatalf("missing order params")
	}
	if len(order.Orderings) != 1 || order.Orderings[0].Column != "title" || order.Orderings[0].Asc {
		t.Fatalf("unexpected order orderings: %+v", order.Orderings)
	}

	if _, ok := ps[filter.SuffixOrder]; ok {
		t.Fatalf("unexpected raw suffix key in param set")
	}
}
