package util_test

import (
	"testing"

	"projectforge.dev/projectforge/app/util"
)

func TestStringToPlural(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"cat", "cats"},
		{"dog", "dogs"},
		{"child", "children"},
		{"goose", "geese"},
		{"man", "men"},
		{"CLASS", "CLASSes"},
		{"BUS", "BUSes"},
		{"PASS", "PASSes"},
		{"ox", "oxen"},
		{"a", "as"},
	}

	for _, tt := range tests {
		x := tt
		result := util.StringToPlural(x.input)
		if result != x.expected {
			t.Errorf("StringToPlural(%s) = %s; expected %s", x.input, result, x.expected)
		}
	}
}

func TestStringToSingular(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"cats", "cat"},
		{"dogs", "dog"},
		{"children", "child"},
		{"geese", "goose"},
		{"men", "man"},
		{"CLASSes", "CLASS"},
		{"BUSes", "BUS"},
		{"PASSes", "PASS"},
		{"oxen", "ox"},
		{"as", "a"},
	}

	for _, tt := range tests {
		x := tt
		result := util.StringToSingular(x.input)
		if result != x.expected {
			t.Errorf("StringToSingular(%s) = %s; expected %s", x.input, result, x.expected)
		}
	}
}

func TestStringForms(t *testing.T) {
	tests := []struct {
		input          string
		expectedSing   string
		expectedPlural string
	}{
		{"cat", "cat", "cats"},
		{"dogs", "dog", "dogs"},
		{"children", "child", "children"},
		{"geese", "goose", "geese"},
		{"CLASS", "CLASS", "CLASSes"},
	}

	for _, tt := range tests {
		x := tt
		singResult, pluralResult := util.StringForms(x.input)
		if singResult != x.expectedSing || pluralResult != x.expectedPlural {
			t.Errorf("StringForms(%s) = (%s, %s); expected (%s, %s)", x.input, singResult, pluralResult, x.expectedSing, x.expectedPlural)
		}
	}
}

func TestStringPlural(t *testing.T) {
	tests := []struct {
		count    int
		input    string
		expected string
	}{
		{1, "cat", "1 cat"},
		{2, "dog", "2 dogs"},
		{0, "child", "0 children"},
		{-1, "goose", "-1 goose"},
		{5, "man", "5 men"},
		{3, "CLASS", "3 CLASSes"},
	}

	for _, tt := range tests {
		x := tt
		result := util.StringPlural(x.count, x.input)
		if result != x.expected {
			t.Errorf("StringPlural(%d, %s) = %s; expected %s", x.count, x.input, result, x.expected)
		}
	}
}
