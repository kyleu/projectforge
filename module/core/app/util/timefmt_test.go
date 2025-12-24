package util_test

import (
	"fmt"
	"testing"
	"time"

	"{{{ .Package }}}/app/util"
)

//nolint:paralleltest,tparallel
func TestTimeFormatting(t *testing.T) {
	t.Parallel()
	testTime := time.Date(2023, 5, 15, 14, 30, 45, 123456789, time.UTC)

	tests := []struct {
		name     string
		function func(*time.Time) string
		expected string
	}{
		{"TimeToFull", util.TimeToFull, "2023-05-15 14:30:45"},
		{"TimeToFullMS", util.TimeToFullMS, "2023-05-15 14:30:45.123456"},
		{"TimeToHours", util.TimeToHours, "14:30:45"},
		{"TimeToHTML", util.TimeToHTML, "2023-05-15T14:30:45"},
		{"TimeToJS", util.TimeToJS, "2023-05-15T14:30:45Z"},
		{"TimeToJSFull", util.TimeToJSFull, "2023-05-15T14:30:45.123Z"},
		{"TimeToRFC3339", util.TimeToRFC3339, "2023-05-15T14:30:45.123Z"},
		{"TimeToVerbose", util.TimeToVerbose, "Mon May 15 14:30:45 2023 +0000"},
		{"TimeToYMD", util.TimeToYMD, "2023-05-15"},
	}

	for _, tt := range tests {
		x := tt
		t.Run(x.name, func(t *testing.T) {
			result := x.function(&testTime)
			if result != x.expected {
				t.Errorf("%s = %v, want %v", x.name, result, x.expected)
			}
		})
	}
}

//nolint:paralleltest,tparallel
func TestTimeParsing(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		function func(string) (*time.Time, error)
		input    string
		expected time.Time
	}{
		{"TimeFromFull", util.TimeFromFull, "2023-05-15 14:30:45", time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)},
		{"TimeFromFullMS", util.TimeFromFullMS, "2023-05-15 14:30:45.123456", time.Date(2023, 5, 15, 14, 30, 45, 123456000, time.UTC)},
		{"TimeFromHTML", util.TimeFromHTML, "2023-05-15T14:30:45", time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)},
		{"TimeFromJS", util.TimeFromJS, "2023-05-15T14:30:45Z", time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)},
		{"TimeFromRFC3339", util.TimeFromRFC3339, "2023-05-15T14:30:45.123Z", time.Date(2023, 5, 15, 14, 30, 45, 123000000, time.UTC)},
		{"TimeFromVerbose", util.TimeFromVerbose, "Mon May 15 14:30:45 2023 +0000", time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)},
		{"TimeFromYMD", util.TimeFromYMD, "2023-05-15", time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		x := tt
		t.Run(x.name, func(t *testing.T) {
			result, err := x.function(x.input)
			if err != nil {
				t.Errorf("%s returned unexpected error: %v", x.name, err)
			}
			if !result.Equal(x.expected) {
				t.Errorf("%s = %v, want %v", x.name, result, x.expected)
			}
		})
	}
}

func TestFormatSeconds(t *testing.T) {
	t.Parallel()
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
		x := tt
		t.Run(fmt.Sprintf("%.3f", x.input), func(t *testing.T) {
			t.Parallel()
			result := util.FormatSeconds(x.input)
			if result != x.expected {
				t.Errorf("FormatSeconds(%.3f) = %s, want %s", x.input, result, x.expected)
			}
		})
	}
}

func TestFormatMilliseconds(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		input           float64
		includeFraction bool
		expected        string
	}{
		{"HourMinuteSecondFraction", 3_723_456, true, "1:02:03.456"},
		{"HourZeroMinutes", 3_603_000, false, "1:00:03"},
		{"MinuteSecondFraction", 123_456, true, "2:03.456"},
		{"SecondsOnlyFraction", 3_456, true, "3.456"},
		{"SecondsOnlyNoFraction", 3_000, false, "3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := util.FormatMilliseconds(tt.input, tt.includeFraction)
			if result != tt.expected {
				t.Errorf("FormatMilliseconds(%f, %t) = %s, want %s", tt.input, tt.includeFraction, result, tt.expected)
			}
		})
	}
}

//nolint:paralleltest,tparallel
func TestFormatSecondsFull(t *testing.T) {
	t.Parallel()
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
			result := util.FormatSecondsFull(tt.input)
			if result != tt.expected {
				t.Errorf("FormatSecondsFull(%.3f) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
