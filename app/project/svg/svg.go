package svg

import (
	"regexp"

	"projectforge.dev/projectforge/app/util"
)

var re = regexp.MustCompile(`\n[ \\t]*`)

type SVG struct {
	Key    string
	Markup string
}

func (s *SVG) Proper() string {
	return util.StringToProper(s.Key)
}

type SVGs []*SVG
