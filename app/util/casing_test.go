//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

func TestStringToInitials(t *testing.T) {
	t.Parallel()
	_ = util.InitAcronyms()

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "simple camelCase", input: "helloWorld", expected: "hw"},
		{name: "PascalCase", input: "HelloWorld", expected: "hw"},
		{name: "snake_case", input: "hello_world", expected: "hw"},
		{name: "kebab-case", input: "hello-world", expected: "hw"},
		{name: "single word", input: "hello", expected: "h"},
		{name: "three words", input: "helloWorldTest", expected: "hwt"},
		{name: "with spaces", input: "hello world", expected: "hw"},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := util.StringToInitials(c.input)
			if result != c.expected {
				t.Errorf("StringToInitials(%s) = %s, expected %s", c.input, result, c.expected)
			}
		})
	}
}

func TestStringToKebab(t *testing.T) {
	t.Parallel()
	_ = util.InitAcronyms()

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "simple camelCase", input: "helloWorld", expected: "hello-world"},
		{name: "PascalCase", input: "HelloWorld", expected: "hello-world"},
		{name: "snake_case", input: "hello_world", expected: "hello-world"},
		{name: "single word", input: "hello", expected: "hello"},
		{name: "with spaces", input: "hello world", expected: "hello-world"},
		{name: "with acronym", input: "parseJSONData", expected: "parse-json-data"},
		{name: "multiple caps", input: "XMLParser", expected: "xml-parser"},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := util.StringToKebab(c.input)
			if result != c.expected {
				t.Errorf("StringToKebab(%s) = %s, expected %s", c.input, result, c.expected)
			}
		})
	}
}

func TestStringCasingEdgeCases(t *testing.T) {
	t.Parallel()
	_ = util.InitAcronyms()

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		if result := util.StringToSnake(""); result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
		if result := util.StringToKebab(""); result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
		if result := util.StringToCamel(""); result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
		if result := util.StringToProper(""); result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
	})

	t.Run("whitespace only", func(t *testing.T) {
		t.Parallel()
		if result := util.StringToSnake("   "); result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
	})

	t.Run("numbers in string", func(t *testing.T) {
		t.Parallel()
		result := util.StringToSnake("test123Value")
		if result != "test_123_value" {
			t.Errorf("expected 'test_123_value', got %s", result)
		}
	})
}
