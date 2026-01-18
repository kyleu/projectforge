//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"{{{ .Package }}}/app/util"
)

func TestConfigureAcronym(t *testing.T) {
	t.Parallel()

	t.Run("configures acronym", func(t *testing.T) {
		t.Parallel()
		util.ConfigureAcronym("TEST", "test")
	})
}

func TestInitAcronyms(t *testing.T) {
	t.Parallel()

	t.Run("initializes default acronyms", func(t *testing.T) {
		t.Parallel()
		_ = util.InitAcronyms()
		result := util.StringToProper("api")
		if result != "API" {
			t.Errorf("expected 'API', got '%s'", result)
		}
	})

	t.Run("handles extra acronyms", func(t *testing.T) {
		t.Parallel()
		_ = util.InitAcronyms()
		result := util.StringToProper("htmlParser", "parser")
		if result != "HTMLParser" {
			t.Errorf("expected 'HTMLParser', got '%s'", result)
		}
	})
}

func TestAcronymProcessing(t *testing.T) {
	t.Parallel()
	_ = util.InitAcronyms()

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "single acronym API", input: "api", expected: "API"},
		{name: "single acronym ID", input: "id", expected: "ID"},
		{name: "single acronym URL", input: "url", expected: "URL"},
		{name: "single acronym HTML", input: "html", expected: "HTML"},
		{name: "single acronym JSON", input: "json", expected: "JSON"},
		{name: "acronym at end", input: "userID", expected: "UserID"},
		{name: "acronym at start", input: "apiEndpoint", expected: "APIEndpoint"},
		{name: "multiple acronyms", input: "htmlToJson", expected: "HTMLToJSON"},
		{name: "no acronym", input: "hello", expected: "Hello"},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := util.StringToProper(c.input)
			if result != c.expected {
				t.Errorf("StringToProper(%s) = %s, expected %s", c.input, result, c.expected)
			}
		})
	}
}
