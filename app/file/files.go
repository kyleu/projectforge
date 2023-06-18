package file

import (
	"github.com/samber/lo"
	"strings"
)

type Files []*File

func (f Files) String() string {
	var sb strings.Builder
	for _, file := range f {
		sb.WriteString(" - ")
		sb.WriteString(file.FullPath())
		sb.WriteString("\n")
	}
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
