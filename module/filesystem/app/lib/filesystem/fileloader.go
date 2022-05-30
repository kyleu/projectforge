package filesystem

import (
	"os"

	"{{{ .Package }}}/app/util"
)

type FileLoader interface {
	Root() string
	Clone() FileLoader
	PeekFile(path string, maxSize int) ([]byte, error)
	ReadFile(path string) ([]byte, error)
	CreateDirectory(path string) error
	WriteFile(path string, content []byte, mode os.FileMode, overwrite bool) error
	CopyFile(src string, tgt string) error
	CopyRecursive(src string, tgt string, ignore []string, logger util.Logger) error
	Move(src string, tgt string) error
	ListFiles(path string, ignore []string, logger util.Logger) []os.DirEntry
	ListFilesRecursive(path string, ignore []string, logger util.Logger) ([]string, error)
	ListJSON(path string, ignore []string, trimExtension bool, logger util.Logger) []string
	ListExtension(path string, ext string, ignore []string, trimExtension bool, logger util.Logger) []string
	ListDirectories(path string, ignore []string, logger util.Logger) []string
	Walk(path string, ign []string, fn func(fp string, info os.FileInfo, err error) error) error
	Stat(path string) (os.FileInfo, error)
	SetMode(path string, mode os.FileMode) error
	Exists(path string) bool
	IsDir(path string) bool
	Remove(path string, logger util.Logger) error
	RemoveRecursive(pt string, logger util.Logger) error
}
