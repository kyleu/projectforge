package file

import (
	"projectforge.dev/projectforge/app/util"
	"strings"

	"github.com/samber/lo"
)

type Files []*File

func (f Files) String() string {
	var sb strings.Builder
	lo.ForEach(f, func(file *File, _ int) {
		sb.WriteString(" - ")
		sb.WriteString(file.FullPath())
		sb.WriteString(util.StringDetectLinebreak(file.Content))
	})
	return sb.String()
}

func (f Files) Get(path string) *File {
	return lo.FindOrElse(f, nil, func(file *File) bool {
		return file.FullPath() == path
	})
}

func (f Files) Exists(path string) bool {
	return f.Get(path) != nil
}
