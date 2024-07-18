package util

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeFormatting(t *testing.T) {
	testTime := time.Date(2023, 5, 15, 14, 30, 45, 123456789, time.UTC)

	tests := []struct {
		name     string
		function func(*time.Time) string
		expected string
	}{
		{"TimeToFull", TimeToFull, "2023-05-15 14:30:45"},
		{"TimeToFullMS", TimeToFullMS, "2023-05-15 14:30:45.123456"},
		{"TimeToHours", TimeToHours, "14:30:45"},
		{"TimeToHTML", TimeToHTML, "2023-05-15T14:30:45"},
		{"TimeToJS", TimeToJS, "2023-05-15T14:30:45Z"},
		{"TimeToJSFull", TimeToJSFull, "2023-05-15T14:30:45.123Z"},
		{"TimeToRFC3339", TimeToRFC3339, "2023-05-15T14:30:45.123Z"},
		{"TimeToVerbose", TimeToVerbose, "Mon May 15 14:30:45 2023 +0000"},
		{"TimeToYMD", TimeToYMD, "2023-05-15"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(&testTime)
			if result != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestTimeParsing(t *testing.T) {
	tests := []struct {
		name     string
		function func(string) (*time.Time, error)
		input    string
		expected time.Time
	}{
		{"TimeFromFull", TimeFromFull, "2023-05-15 14:30:45", time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)},
		{"TimeFromFullMS", TimeFromFullMS, "2023-05-15 14:30:45.123456", time.Date(2023, 5, 15, 14, 30, 45, 123456000, time.UTC)},
		{"TimeFromHTML", TimeFromHTML, "2023-05-15T14:30:45", time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)},
		{"TimeFromJS", TimeFromJS, "2023-05-15T14:30:45Z", time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)},
		{"TimeFromRFC3339", TimeFromRFC3339, "2023-05-15T14:30:45.123Z", time.Date(2023, 5, 15, 14, 30, 45, 123000000, time.UTC)},
		{"TimeFromVerbose", TimeFromVerbose, "Mon May 15 14:30:45 2023 +0000", time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)},
		{"TimeFromYMD", TimeFromYMD, "2023-05-15", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.function(tt.input)
			if err != nil {
				t.Errorf("%s returned unexpected error: %v", tt.name, err)
			}
			if !result.Equal(tt.expected) {
				t.Errorf("%s = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestFormatSeconds(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{3661.123, "1h 1m 1.123s"},
		{60.5, "1m 0.500s"},
		{3.0, "3s"},
		{0.123, "0.123s"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.3f", tt.input), func(t *testing.T) {
			result := FormatSeconds(tt.input)
			if result != tt.expected {
				t.Errorf("FormatSeconds(%.3f) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatSecondsFull(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{3661.123, "1 hour, 1 minute, 1 second, 123 milliseconds"},
		{60.5, "1 minute, 0 seconds, 500 milliseconds"},
		{3.0, "3 seconds"},
		{0.123, "0 seconds, 123 milliseconds"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%.3f", tt.input), func(t *testing.T) {
			result := FormatSecondsFull(tt.input)
			if result != tt.expected {
				t.Errorf("FormatSecondsFull(%.3f) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
