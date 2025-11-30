package util

import (
	"fmt"
	"math"
)

func MicrosToMillis(i int) string {
	div := 1000

	ms := i / div
	if ms >= 20 {
		return fmt.Sprintf("%dms", ms)
	}

	x := float64(ms) + (float64(i%div) / float64(div))
	if x == math.Round(x) {
		return fmt.Sprintf("%dms", ms)
	}

	return fmt.Sprintf("%.3f", x) + "ms"
}

func MillisToString(i int) string {
	if i == 0 {
		return "0s"
	}

	negative := i < 0
	if negative {
		i = -i
	}

	seconds := i / 1000
	if seconds == 0 {
		return "0s"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	seconds = seconds % 60

	var ret string
	if hours > 0 {
		ret = fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	} else {
		ret = fmt.Sprintf("%d:%02d", minutes, seconds)
	}
	if negative {
		return "-" + ret
	}
	return ret
}

func Percent(f float64) string {
	return fmt.Sprintf("%.3f", f) + "%"
}
