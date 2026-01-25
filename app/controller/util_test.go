package controller_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"

	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/csession"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/util"
)

func TestERsp(t *testing.T) {
	t.Parallel()
	msg, err := controller.ERsp("boom %d", 1)
	if err == nil || !strings.Contains(err.Error(), "boom 1") {
		t.Fatalf("ERsp error was [%v]", err)
	}
	if msg != "" {
		t.Fatalf("ERsp msg was [%v]", msg)
	}
}

func TestFlashAndRedir(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	wc := cutil.NewWriteCounter(rr)
	ps := &cutil.PageState{
		Session: util.ValueMap{},
		W:       wc,
		Logger:  zap.NewNop().Sugar(),
	}

	redir, err := controller.FlashAndRedir(true, "ok", "/next", ps)
	if err != nil {
		t.Fatalf("FlashAndRedir returned error: %v", err)
	}
	if redir != "/next" {
		t.Fatalf("FlashAndRedir returned [%v]", redir)
	}
	if got := ps.Session.GetStringOpt(csession.WebFlashKey); got != "success:ok" {
		t.Fatalf("flash session was [%v]", got)
	}
	if len(rr.Result().Cookies()) == 0 {
		t.Fatalf("expected session cookie to be set")
	}

	redir, err = controller.FlashAndRedir(false, "bad", "/again", ps)
	if err != nil {
		t.Fatalf("FlashAndRedir returned error: %v", err)
	}
	if redir != "/again" {
		t.Fatalf("FlashAndRedir returned [%v]", redir)
	}
	if got := ps.Session.GetStringOpt(csession.WebFlashKey); got != "success:ok;error:bad" {
		t.Fatalf("flash session was [%v]", got)
	}
}

func TestFlashAndRedirRejectsExternal(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	wc := cutil.NewWriteCounter(rr)
	ps := &cutil.PageState{
		Session: util.ValueMap{},
		W:       wc,
		Logger:  zap.NewNop().Sugar(),
	}

	redir, err := controller.FlashAndRedir(true, "ok", "http://example.com", ps)
	if err != nil {
		t.Fatalf("FlashAndRedir returned error: %v", err)
	}
	if redir != "/" {
		t.Fatalf("FlashAndRedir returned [%v]", redir)
	}
}

func TestReturnToReferrer(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	wc := cutil.NewWriteCounter(rr)
	ps := &cutil.PageState{
		Session: util.ValueMap{csession.ReferKey: "/back"},
		W:       wc,
		Logger:  zap.NewNop().Sugar(),
	}

	redir, err := controller.ReturnToReferrer("done", "/default", ps)
	if err != nil {
		t.Fatalf("ReturnToReferrer returned error: %v", err)
	}
	if redir != "/back" {
		t.Fatalf("ReturnToReferrer returned [%v]", redir)
	}
	if _, ok := ps.Session[csession.ReferKey]; ok {
		t.Fatalf("refer key was not removed")
	}
}

func TestOptionsAndHead(t *testing.T) {
	t.Parallel()
	rr := httptest.NewRecorder()
	controller.Options(rr, httptest.NewRequest(http.MethodOptions, "http://example.com", http.NoBody))
	if rr.Code != http.StatusOK {
		t.Fatalf("Options status was [%d]", rr.Code)
	}
	if rr.Header().Get(cutil.HeaderAccessControlAllowOrigin) == "" {
		t.Fatalf("Options missing CORS headers")
	}

	rr = httptest.NewRecorder()
	controller.Head(rr, httptest.NewRequest(http.MethodHead, "http://example.com", http.NoBody))
	if rr.Code != http.StatusOK {
		t.Fatalf("Head status was [%d]", rr.Code)
	}
	if rr.Header().Get(cutil.HeaderAccessControlAllowOrigin) == "" {
		t.Fatalf("Head missing CORS headers")
	}
}
