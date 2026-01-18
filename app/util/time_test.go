//go:build test_all || !func_test
// +build test_all !func_test

package util_test

import (
	"testing"
	"time"

	"projectforge.dev/projectforge/app/util"
)

const neverStr = "<never>"

func TestTimeTruncate(t *testing.T) {
	t.Parallel()

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		result := util.TimeTruncate(nil)
		if result != nil {
			t.Error("expected nil")
		}
	})

	t.Run("truncates time to midnight", func(t *testing.T) {
		t.Parallel()
		tm := time.Date(2024, 6, 15, 14, 30, 45, 123456789, time.UTC)
		result := util.TimeTruncate(&tm)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 || result.Nanosecond() != 0 {
			t.Errorf("expected midnight, got %v", result)
		}
		if result.Year() != 2024 || result.Month() != 6 || result.Day() != 15 {
			t.Errorf("date changed unexpectedly: %v", result)
		}
	})
}

func TestTimeCurrent(t *testing.T) {
	t.Parallel()

	t.Run("returns current time", func(t *testing.T) {
		t.Parallel()
		before := time.Now()
		result := util.TimeCurrent()
		after := time.Now()
		if result.Before(before) || result.After(after) {
			t.Error("TimeCurrent returned time outside expected range")
		}
	})
}

func TestTimeCurrentP(t *testing.T) {
	t.Parallel()

	t.Run("returns non-nil pointer", func(t *testing.T) {
		t.Parallel()
		result := util.TimeCurrentP()
		if result == nil {
			t.Error("expected non-nil pointer")
		}
	})
}

func TestTimeToday(t *testing.T) {
	t.Parallel()

	t.Run("returns today at midnight", func(t *testing.T) {
		t.Parallel()
		result := util.TimeToday()
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
			t.Errorf("expected midnight, got %v", result)
		}
		now := time.Now()
		if result.Year() != now.Year() || result.Month() != now.Month() || result.Day() != now.Day() {
			t.Errorf("expected today's date, got %v", result)
		}
	})
}

func TestTimeCurrentUnix(t *testing.T) {
	t.Parallel()

	t.Run("returns unix timestamp", func(t *testing.T) {
		t.Parallel()
		before := time.Now().Unix()
		result := util.TimeCurrentUnix()
		after := time.Now().Unix()
		if result < before || result > after {
			t.Error("TimeCurrentUnix returned timestamp outside expected range")
		}
	})
}

func TestTimeCurrentMillis(t *testing.T) {
	t.Parallel()

	t.Run("returns milliseconds", func(t *testing.T) {
		t.Parallel()
		before := time.Now().UnixMilli()
		result := util.TimeCurrentMillis()
		after := time.Now().UnixMilli()
		if result < before || result > after {
			t.Error("TimeCurrentMillis returned timestamp outside expected range")
		}
	})
}

func TestTimeCurrentNanos(t *testing.T) {
	t.Parallel()

	t.Run("returns nanoseconds", func(t *testing.T) {
		t.Parallel()
		before := time.Now().UnixNano()
		result := util.TimeCurrentNanos()
		after := time.Now().UnixNano()
		if result < before || result > after {
			t.Error("TimeCurrentNanos returned timestamp outside expected range")
		}
	})
}

func TestTimeRelative(t *testing.T) {
	t.Parallel()

	t.Run("nil returns <never>", func(t *testing.T) {
		t.Parallel()
		result := util.TimeRelative(nil)
		if result != neverStr {
			t.Errorf("expected '%s', got %s", neverStr, result)
		}
	})

	t.Run("recent time returns relative string", func(t *testing.T) {
		t.Parallel()
		tm := time.Now().Add(-5 * time.Minute)
		result := util.TimeRelative(&tm)
		if result == "" || result == neverStr {
			t.Errorf("expected relative time string, got %s", result)
		}
	})
}

func TestTimeRoundedP(t *testing.T) {
	t.Parallel()

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		result := util.TimeRoundedP(nil, time.Hour)
		if result != nil {
			t.Error("expected nil")
		}
	})

	t.Run("rounds to nearest duration", func(t *testing.T) {
		t.Parallel()
		tm := time.Date(2024, 6, 15, 14, 37, 45, 0, time.UTC)
		result := util.TimeRoundedP(&tm, time.Hour)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Minute() != 0 || result.Second() != 0 {
			t.Errorf("expected rounded to hour, got %v", result)
		}
	})
}

func TestTimeRounded(t *testing.T) {
	t.Parallel()

	t.Run("zero time returns zero", func(t *testing.T) {
		t.Parallel()
		var zero time.Time
		result := util.TimeRounded(zero, time.Hour)
		if !result.IsZero() {
			t.Error("expected zero time")
		}
	})

	t.Run("rounds to nearest duration", func(t *testing.T) {
		t.Parallel()
		tm := time.Date(2024, 6, 15, 14, 37, 45, 0, time.UTC)
		result := util.TimeRounded(tm, time.Hour)
		if result.Minute() != 0 || result.Second() != 0 {
			t.Errorf("expected rounded to hour, got %v", result)
		}
	})
}

func TestTimeToMap(t *testing.T) {
	t.Parallel()

	t.Run("returns map with epoch and iso8601", func(t *testing.T) {
		t.Parallel()
		tm := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
		result := util.TimeToMap(tm)
		if _, ok := result["epoch"]; !ok {
			t.Error("expected epoch key")
		}
		if _, ok := result["iso8601"]; !ok {
			t.Error("expected iso8601 key")
		}
		epoch, ok := result["epoch"].(int64)
		if !ok {
			t.Error("epoch is not int64")
		}
		if epoch != tm.UnixMilli() {
			t.Errorf("expected epoch %d, got %v", tm.UnixMilli(), result["epoch"])
		}
	})
}

func TestTimeToString(t *testing.T) {
	t.Parallel()

	t.Run("nil returns empty string", func(t *testing.T) {
		t.Parallel()
		result := util.TimeToString(nil, "2006-01-02")
		if result != "" {
			t.Errorf("expected empty string, got %s", result)
		}
	})

	t.Run("formats time correctly", func(t *testing.T) {
		t.Parallel()
		tm := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
		result := util.TimeToString(&tm, "2006-01-02")
		if result != "2024-06-15" {
			t.Errorf("expected '2024-06-15', got %s", result)
		}
	})
}

func TestTimeFromString(t *testing.T) {
	t.Parallel()

	t.Run("empty string returns nil", func(t *testing.T) {
		t.Parallel()
		result, err := util.TimeFromString("")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != nil {
			t.Error("expected nil")
		}
	})

	t.Run("valid date string parses correctly", func(t *testing.T) {
		t.Parallel()
		result, err := util.TimeFromString("2024-06-15")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Year() != 2024 || result.Month() != 6 || result.Day() != 15 {
			t.Errorf("unexpected date: %v", result)
		}
	})

	t.Run("invalid string returns error", func(t *testing.T) {
		t.Parallel()
		_, err := util.TimeFromString("not-a-date")
		if err == nil {
			t.Error("expected error for invalid date string")
		}
	})
}

func TestTimeFromStringSimple(t *testing.T) {
	t.Parallel()

	t.Run("valid date string parses", func(t *testing.T) {
		t.Parallel()
		result := util.TimeFromStringSimple("2024-06-15")
		if result == nil {
			t.Error("expected non-nil result")
		}
	})

	t.Run("invalid string returns nil", func(t *testing.T) {
		t.Parallel()
		result := util.TimeFromStringSimple("not-a-date")
		if result != nil {
			t.Error("expected nil for invalid date")
		}
	})
}

func TestTimePlusDays(t *testing.T) {
	t.Parallel()

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		result := util.TimePlusDays(nil, 5)
		if result != nil {
			t.Error("expected nil")
		}
	})

	t.Run("adds days correctly", func(t *testing.T) {
		t.Parallel()
		tm := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
		result := util.TimePlusDays(&tm, 5)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Day() != 20 {
			t.Errorf("expected day 20, got %d", result.Day())
		}
	})

	t.Run("subtracts days with negative value", func(t *testing.T) {
		t.Parallel()
		tm := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
		result := util.TimePlusDays(&tm, -5)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Day() != 10 {
			t.Errorf("expected day 10, got %d", result.Day())
		}
	})
}

func TestTimeMax(t *testing.T) {
	t.Parallel()

	t.Run("all nil returns nil", func(t *testing.T) {
		t.Parallel()
		result := util.TimeMax(nil, nil, nil)
		if result != nil {
			t.Error("expected nil")
		}
	})

	t.Run("returns max time", func(t *testing.T) {
		t.Parallel()
		t1 := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
		t2 := time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC)
		t3 := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
		result := util.TimeMax(&t1, &t2, &t3)
		if result == nil || !result.Equal(t2) {
			t.Errorf("expected t2, got %v", result)
		}
	})

	t.Run("ignores nil values", func(t *testing.T) {
		t.Parallel()
		t1 := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
		result := util.TimeMax(nil, &t1, nil)
		if result == nil || !result.Equal(t1) {
			t.Errorf("expected t1, got %v", result)
		}
	})
}

func TestTimeMin(t *testing.T) {
	t.Parallel()

	t.Run("all nil returns nil", func(t *testing.T) {
		t.Parallel()
		result := util.TimeMin(nil, nil, nil)
		if result != nil {
			t.Error("expected nil")
		}
	})

	t.Run("returns min time", func(t *testing.T) {
		t.Parallel()
		t1 := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
		t2 := time.Date(2024, 6, 20, 0, 0, 0, 0, time.UTC)
		t3 := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
		result := util.TimeMin(&t1, &t2, &t3)
		if result == nil || !result.Equal(t1) {
			t.Errorf("expected t1, got %v", result)
		}
	})

	t.Run("ignores nil values", func(t *testing.T) {
		t.Parallel()
		t1 := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
		result := util.TimeMin(nil, &t1, nil)
		if result == nil || !result.Equal(t1) {
			t.Errorf("expected t1, got %v", result)
		}
	})
}
