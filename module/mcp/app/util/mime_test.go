//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"{{{ .Package }}}/app/util"
)

func TestMIMEHelpers(t *testing.T) {
	if ext := util.ExtensionFromMIME(""); ext != "json" {
		t.Fatalf("expected json extension, got %q", ext)
	}
	if ext := util.ExtensionFromMIME(util.MIMETypeMarkdown); ext != "md" {
		t.Fatalf("expected md extension, got %q", ext)
	}
	if mt := util.MIMEFromExtension("md"); mt != util.MIMETypeMarkdown {
		t.Fatalf("expected markdown mime, got %q", mt)
	}
	if mt := util.MIMEFromExtension("unknownext"); mt != util.MIMETypeJSON {
		t.Fatalf("expected json mime for unknown ext, got %q", mt)
	}
}
