// Content managed by Project Forge, see [projectforge.md] for details.
package util

import (
	"time"
)

func TimerStart() int64 {
	return time.Now().UnixNano()
}

func TimerEnd(startNanos int64) int {
	return int((time.Now().UnixNano() - startNanos) / int64(time.Microsecond))
}
