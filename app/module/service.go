package module

import (
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Service struct {
	filesystems map[string]filesystem.FileLoader
	logger      *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger) *Service {
	return &Service{filesystems: map[string]filesystem.FileLoader{}, logger: logger}
}

func (s *Service) GetFilesystem(mod *Module) filesystem.FileLoader {
	curr, ok := s.filesystems[mod.Key]
	if ok {
		return curr
	}

	p := filepath.Join("module", mod.Key)
	fs := filesystem.NewFileSystem(p, s.logger)

	s.filesystems[mod.Key] = fs
	return fs
}

func (s *Service) GetFiles(mod *Module, changes *file.Changeset, addHeader bool, tgt filesystem.FileLoader) (file.Files, error) {
	loader := s.GetFilesystem(mod)
	fs, err := loader.ListFilesRecursive("", nil)
	if err != nil {
		return nil, err
	}
	ret := make(file.Files, 0, len(fs))
	for _, f := range fs {
		if f == ".module.json" {
			continue
		}
		f = strings.TrimPrefix(strings.TrimPrefix(f, mod.Path()), "/")
		mode, b, err := mod.FileContent(loader, f)
		if err != nil {
			return nil, err
		}
		fl := file.NewFile(f, mode, b, addHeader, s.logger)
		if strings.Contains(fl.Content, file.SectionPrefix) && tgt.Exists(f) {
			tgtBytes, _ := tgt.ReadFile(f)
			if utf8.Valid(tgtBytes) {
				newContent, err := file.CopySections(string(tgtBytes), fl.Content)
				if err != nil {
					return nil, errors.Wrapf(err, "error reading sections from [%s]", f)
				}
				fl.Content = newContent
			}
		}
		fl = fl.Apply(changes)
		ret = append(ret, fl)
	}
	return ret, nil
}

func (s *Service) GetModules(keys ...string) (Modules, error) {
	ret := Modules{}
	for _, m := range keys {
		mod, ok := AvailableModules[m]
		if !ok {
			return nil, errors.New("no module available with key [" + m + "]")
		}
		ret[m] = mod
	}
	return ret, nil
}
