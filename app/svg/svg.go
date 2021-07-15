package svg

import (
	"regexp"
)

var re = regexp.MustCompile(`\n[ \\t]*`)

type SVG struct {
	Key    string
	Markup string
}
