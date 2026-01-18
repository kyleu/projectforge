//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"

	"{{{ .Package }}}/app/util"
)

func TestByteSizeSI(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    int64
		expected string
	}{
		{name: "zero bytes", input: 0, expected: "0 B"},
		{name: "small bytes", input: 500, expected: "500 B"},
		{name: "exactly 1000 bytes", input: 1000, expected: "1.0 KB"},
		{name: "kilobytes", input: 1500, expected: "1.5 KB"},
		{name: "megabytes", input: 1500000, expected: "1.5 MB"},
		{name: "gigabytes", input: 1500000000, expected: "1.5 GB"},
		{name: "terabytes", input: 1500000000000, expected: "1.5 TB"},
		{name: "petabytes", input: 1500000000000000, expected: "1.5 PB"},
		{name: "exact kilobyte", input: 1000, expected: "1.0 KB"},
		{name: "exact megabyte", input: 1000000, expected: "1.0 MB"},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := util.ByteSizeSI(c.input)
			if result != c.expected {
				t.Errorf("ByteSizeSI(%d) = %s, expected %s", c.input, result, c.expected)
			}
		})
	}
}

func TestByteSizeIEC(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    int64
		expected string
	}{
		{name: "zero bytes", input: 0, expected: "0 B"},
		{name: "small bytes", input: 500, expected: "500 B"},
		{name: "exactly 1024 bytes", input: 1024, expected: "1.0 KiB"},
		{name: "kibibytes", input: 1536, expected: "1.5 KiB"},
		{name: "mebibytes", input: 1572864, expected: "1.5 MiB"},
		{name: "gibibytes", input: 1610612736, expected: "1.5 GiB"},
		{name: "tebibytes", input: 1649267441664, expected: "1.5 TiB"},
		{name: "exact kibibyte", input: 1024, expected: "1.0 KiB"},
		{name: "exact mebibyte", input: 1048576, expected: "1.0 MiB"},
	}

	for _, tc := range cases {
		c := tc
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := util.ByteSizeIEC(c.input)
			if result != c.expected {
				t.Errorf("ByteSizeIEC(%d) = %s, expected %s", c.input, result, c.expected)
			}
		})
	}
}
