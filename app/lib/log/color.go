package log

import (
	"fmt"
)

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

type Color uint8

func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

var levelToColor = map[string]Color{"debug": Magenta, "info": Cyan, "warn": Yellow, "error": Red, "dpanic": Red, "panic": Red, "fatal": Red}
