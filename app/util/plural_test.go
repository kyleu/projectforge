package util

import (
	"testing"
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

	for _, test := range tests {
		result := StringToPlural(test.input)
		if result != test.expected {
			t.Errorf("StringToPlural(%s) = %s; expected %s", test.input, result, test.expected)
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

	for _, test := range tests {
		result := StringToSingular(test.input)
		if result != test.expected {
			t.Errorf("StringToSingular(%s) = %s; expected %s", test.input, result, test.expected)
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

	for _, test := range tests {
		singResult, pluralResult := StringForms(test.input)
		if singResult != test.expectedSing || pluralResult != test.expectedPlural {
			t.Errorf("StringForms(%s) = (%s, %s); expected (%s, %s)", test.input, singResult, pluralResult, test.expectedSing, test.expectedPlural)
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

	for _, test := range tests {
		result := StringPlural(test.count, test.input)
		if result != test.expected {
			t.Errorf("StringPlural(%d, %s) = %s; expected %s", test.count, test.input, result, test.expected)
		}
	}
}
