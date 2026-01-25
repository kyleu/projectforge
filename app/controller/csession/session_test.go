package csession_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"

	"projectforge.dev/projectforge/app/controller/csession"
	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

const testCookieValue = "value"

func TestNewCookie(t *testing.T) {
	t.Parallel()
	c := csession.NewCookie(testCookieValue)
	if c.Name != util.AppKey {
		t.Fatalf("cookie name was [%v]", c.Name)
	}
	if c.Value != testCookieValue {
		t.Fatalf("cookie value was [%v]", c.Value)
	}
	if c.Path != "/" {
		t.Fatalf("cookie path was [%v]", c.Path)
	}
	if c.MaxAge <= 0 {
		t.Fatalf("cookie MaxAge was [%v]", c.MaxAge)
	}
	if !c.HttpOnly {
		t.Fatalf("cookie HttpOnly was false")
	}
	if c.SameSite != http.SameSiteLaxMode {
		t.Fatalf("cookie SameSite was [%v]", c.SameSite)
	}
}

func TestGetFromSession(t *testing.T) {
	t.Parallel()
	sess := util.ValueMap{"ok": testCookieValue, "bad": 12}

	if got, err := csession.GetFromSession("ok", sess); err != nil || got != testCookieValue {
		t.Fatalf("GetFromSession returned [%v], err=[%v]", got, err)
	}

	if _, err := csession.GetFromSession("missing", sess); err == nil {
		t.Fatalf("expected error for missing key")
	}

	if _, err := csession.GetFromSession("bad", sess); err == nil {
		t.Fatalf("expected error for non-string value")
	}
}

func TestSaveProfileDefaultRemoves(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	sess := util.ValueMap{"profile": "{\"name\":\"Bob\"}"}

	if err := csession.SaveProfile(user.DefaultProfile.Clone(), rr, sess, zap.NewNop().Sugar()); err != nil {
		t.Fatalf("SaveProfile returned error: %v", err)
	}
	if _, ok := sess["profile"]; ok {
		t.Fatalf("expected profile to be removed")
	}
	if len(rr.Result().Cookies()) == 0 {
		t.Fatalf("expected session cookie to be set")
	}
}

func TestSaveProfileCustomStores(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	sess := util.ValueMap{}

	prof := &user.Profile{Name: "Alice", Mode: "dark"}
	if err := csession.SaveProfile(prof, rr, sess, zap.NewNop().Sugar()); err != nil {
		t.Fatalf("SaveProfile returned error: %v", err)
	}
	if got, ok := sess["profile"]; !ok {
		t.Fatalf("expected profile to be stored")
	} else if got != util.ToJSON(prof) {
		t.Fatalf("profile stored as [%v]", got)
	}
	if len(rr.Result().Cookies()) == 0 {
		t.Fatalf("expected session cookie to be set")
	}
}

func TestSaveProfileNilRemoves(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	sess := util.ValueMap{"profile": "{\"name\":\"Bob\"}"}

	if err := csession.SaveProfile(nil, rr, sess, zap.NewNop().Sugar()); err != nil {
		t.Fatalf("SaveProfile returned error: %v", err)
	}
	if _, ok := sess["profile"]; ok {
		t.Fatalf("expected profile to be removed")
	}
}
