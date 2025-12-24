package util_test

import (
	"fmt"
	"testing"

	"{{{ .Package }}}/app/util"
)

func TestMillisToString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0s"},
		{999, "0s"},
		{1_000, "0:01"},
		{10_000, "0:10"},
		{100_000, "1:40"},
		{3_601_000, "1:00:01"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.input), func(t *testing.T) {
			t.Parallel()
			result := util.MillisToString(tt.input)
			if result != tt.expected {
				t.Errorf("MillisToString(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMicrosToMillis(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0ms"},
		{1, "0.001ms"},
		{1_000, "1ms"},
		{19_000, "19ms"},
		{19_999, "19.999ms"},
		{20_000, "20ms"},
		{20_500, "20ms"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.input), func(t *testing.T) {
			t.Parallel()
			result := util.MicrosToMillis(tt.input)
			if result != tt.expected {
				t.Errorf("MicrosToMillis(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
