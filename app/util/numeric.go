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

func Percent(f float64) string {
	return fmt.Sprintf("%.3f", f) + "%"
}
