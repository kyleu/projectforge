package module

import (
	"os"
	"strings"
	"unicode/utf8"

	"github.com/kyleu/projectforge/app/file"
	"github.com/kyleu/projectforge/app/filesystem"
	"github.com/kyleu/projectforge/app/util"
	"github.com/pkg/errors"
)

const configFilename = ".module.json"

func fileContent(files filesystem.FileLoader, path string) (os.FileMode, []byte, error) {
	stat, err := files.Stat(path)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "file [%s] not found", path)
	}
	b, err := files.ReadFile(path)
	if err != nil {
		return 0, nil, err
	}
	return stat.Mode(), b, nil
}

func (s *Service) LoadAll(keys ...string) (Modules, error) {
	var ret Modules
	for _, key := range keys {
		m, err := s.Load(key)
		if err != nil {
			return nil, err
		}
		ret = append(ret, m)
	}
	return ret, nil
}

func (s *Service) Load(key string) (*Module, error) {
	fs := s.GetFilesystem(key)
	if !fs.Exists(configFilename) {
		msg := "file [%s] does not exist in path for module [%s]"
		return nil, errors.Errorf(msg, configFilename, key)
	}

	b, err := fs.ReadFile(configFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read [%s]", configFilename)
	}

	ret := &Module{}
	err = util.FromJSON(b, ret)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse [%s] as module with key [%s]", configFilename, key)
	}
	ret.Key = key
	s.cache[key] = ret
	return ret, nil
}

func (s *Service) GetFiles(mods Modules, changes *file.Changeset, addHeader bool, tgt filesystem.FileLoader) (file.Files, error) {
	loader := s.getNestedFilesystem(mods)
	fs, err := loader.ListFilesRecursive("", nil)
	if err != nil {
		return nil, err
	}
	ret := make(file.Files, 0, len(fs))
	for _, f := range fs {
		if f == configFilename {
			continue
		}
		// f = strings.TrimPrefix(strings.TrimPrefix(f, mod.Path()), "/")
		mode, b, err := fileContent(loader, f)
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
