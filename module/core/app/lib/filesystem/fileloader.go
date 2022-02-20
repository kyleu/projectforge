package filesystem

import (
	"os"
)

type FileLoader interface {
	Root() string
	Clone() FileLoader
	PeekFile(path string, maxSize int) ([]byte, error)
	ReadFile(path string) ([]byte, error)
	CreateDirectory(path string) error
	WriteFile(path string, content []byte, mode os.FileMode, overwrite bool) error
	CopyFile(src string, tgt string) error
	CopyRecursive(src string, tgt string, ignore []string) error
	Move(src string, tgt string) error
	ListFiles(path string, ignore []string) []os.DirEntry
	ListFilesRecursive(path string, ignore []string) ([]string, error)
	ListJSON(path string, trimExtension bool) []string
	ListExtension(path string, ext string, trimExtension bool) []string
	ListDirectories(path string) []string
	Walk(path string, ign []string, fn func(fp string, info os.FileInfo, err error) error) error
	Stat(path string) (os.FileInfo, error)
	Exists(path string) bool
	IsDir(path string) bool
	Remove(path string) error
	RemoveRecursive(pt string) error
}
