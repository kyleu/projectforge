package cutil_test

import (
	"strings"
	"testing"

	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/user"
	"{{{ .Package }}}/app/util"
)

const pageTitle = "Hello"

func TestPageStateTitleAndData(t *testing.T) {
	t.Parallel()
	ps := &cutil.PageState{}
	ps.SetTitleAndData(pageTitle, 123)

	if ps.Title != pageTitle || ps.Data != 123 {
		t.Fatalf("SetTitleAndData set [%v], [%v]", ps.Title, ps.Data)
	}
}

func TestPageStateTitleString(t *testing.T) {
	t.Parallel()
	ps := &cutil.PageState{}
	if got := ps.TitleString(); got != util.AppName {
		t.Fatalf("TitleString was [%v]", got)
	}

	ps.Title = "!Raw"
	if got := ps.TitleString(); got != "Raw" {
		t.Fatalf("TitleString was [%v]", got)
	}

	ps.Title = pageTitle
	if got := ps.TitleString(); got != pageTitle+" - "+util.AppName {
		t.Fatalf("TitleString was [%v]", got)
	}
}

func TestPageStateAddIcon(t *testing.T) {
	t.Parallel()
	ps := &cutil.PageState{}
	ps.AddIcon("a", "b", "a")

	if len(ps.Icons) != 2 {
		t.Fatalf("expected 2 icons, got %d", len(ps.Icons))
	}
}

func TestPageStateClassDecl(t *testing.T) {
	t.Parallel()
	ps := &cutil.PageState{Icons: []string{"icon"}, Profile: &user.Profile{Mode: "dark"}}
	ps.Browser = "firefox"
	ps.OS = "linux"
	ps.Platform = "desktop"

	if got := ps.ClassDecl(); got != " class=\"mode-dark browser-firefox os-linux platform-desktop\"" {
		t.Fatalf("ClassDecl was [%v]", got)
	}

	ps = &cutil.PageState{}
	if got := ps.ClassDecl(); got != "-" {
		t.Fatalf("ClassDecl was [%v]", got)
	}
}

func TestPageStateMainClasses(t *testing.T) {
	t.Parallel()
	ps := &cutil.PageState{HideHeader: true, HideMenu: true}
	if got := ps.MainClasses(); got != "noheader nomenu" {
		t.Fatalf("MainClasses was [%v]", got)
	}
}

func TestPageStateExtra(t *testing.T) {
	t.Parallel()
	ps := &cutil.PageState{ExtraContent: map[string]string{"x": "y"}}
	if got := ps.Extra("x"); got != "y" {
		t.Fatalf("Extra was [%v]", got)
	}
	if got := ps.Extra("missing"); got != "" {
		t.Fatalf("Extra was [%v]", got)
	}
}

func TestPageStateAddHeaderScript(t *testing.T) {
	t.Parallel()
	ps := &cutil.PageState{}
	ps.AddHeaderScript("client.js", true)

	if !strings.Contains(ps.HeaderContent, "script") || !strings.Contains(ps.HeaderContent, "client.js") {
		t.Fatalf("HeaderContent was [%v]", ps.HeaderContent)
	}
}

func TestPageStateUsername(t *testing.T) {
	t.Parallel()
	ps := &cutil.PageState{Profile: &user.Profile{Name: "alice"}}
	if got := ps.Username(); got != "alice" {
		t.Fatalf("Username was [%v]", got)
	}
}
