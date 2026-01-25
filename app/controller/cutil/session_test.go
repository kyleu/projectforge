package cutil_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"projectforge.dev/projectforge/app/controller/csession"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

const requestBody = "hello"

func TestLoadPageStateNoCookie(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	wc := cutil.NewWriteCounter(rr)
	req := httptest.NewRequest(http.MethodGet, "http://example.com/path", http.NoBody)

	ps := cutil.LoadPageState(nil, wc, req, "test", zap.NewNop().Sugar())
	if ps.Profile == nil || ps.Profile.Name != user.DefaultProfile.Name {
		t.Fatalf("profile was [%v]", ps.Profile)
	}
	if len(ps.Session) != 0 {
		t.Fatalf("session was [%v]", ps.Session)
	}
	if len(ps.Flashes) != 0 {
		t.Fatalf("flashes were [%v]", ps.Flashes)
	}
}

func TestLoadPageStateWithCookieAndFlashes(t *testing.T) {
	t.Parallel()
	logger := zap.NewNop().Sugar()
	sess := util.ValueMap{
		csession.WebFlashKey: "success:ok;error:bad",
		"profile":            util.ToJSON(&user.Profile{Name: "", Mode: "dark"}),
	}
	enc, err := util.EncryptMessage(nil, util.ToJSONCompact(sess), logger)
	if err != nil {
		t.Fatalf("EncryptMessage returned error: %v", err)
	}

	rr := httptest.NewRecorder()
	wc := cutil.NewWriteCounter(rr)
	req := httptest.NewRequest(http.MethodPost, "http://example.com/path", bytes.NewBufferString(requestBody))
	req.AddCookie(csession.NewCookie(enc))

	ps := cutil.LoadPageState(nil, wc, req, "test", logger)
	if len(ps.Flashes) != 2 {
		t.Fatalf("flashes were [%v]", ps.Flashes)
	}
	if _, ok := ps.Session[csession.WebFlashKey]; ok {
		t.Fatalf("expected flash key to be removed from session")
	}
	if ps.Profile == nil || ps.Profile.Name != user.DefaultProfile.Name || ps.Profile.Mode != "dark" {
		t.Fatalf("profile was [%v]", ps.Profile)
	}
	if string(ps.RequestBody) != requestBody {
		t.Fatalf("RequestBody was [%s]", string(ps.RequestBody))
	}
	if len(rr.Result().Cookies()) == 0 {
		t.Fatalf("expected session cookie to be set after flash removal")
	}
}
