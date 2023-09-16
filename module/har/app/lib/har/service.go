package har

import (
	"context"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/lib/search/result"
	"{{{ .Package }}}/app/util"
)

type Service struct {
	FS filesystem.FileLoader
}

func NewService(fs filesystem.FileLoader) *Service {
	return &Service{FS: fs}
}

func (s *Service) List(logger util.Logger) []string {
	return s.FS.ListExtension("./har", "har", nil, true, logger)
}

func (s *Service) Load(fn string) (*Log, error) {
	key := fn
	if !strings.HasSuffix(fn, Ext) {
		fn += Ext
	}
	if !strings.Contains(fn, "har/") {
		fn = path.Join("har", fn)
	}
	if !s.FS.Exists(fn) {
		return nil, errors.Errorf("missing file [%s]", fn)
	}
	b, err := s.FS.ReadFile(fn)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading file [%s]", fn)
	}
	ret := &Wrapper{}
	err = util.FromJSON(b, ret)
	if err != nil {
		return nil, errors.Wrapf(err, "error decoding file [%s]", fn)
	}
	ret.Log.Key = strings.TrimSuffix(key, Ext)
	ret.Log.Entries = ret.Log.Entries.Trimmed().WithJSON()
	return ret.Log, nil
}

func (s *Service) Delete(key string, logger util.Logger) error {
	fn := key
	if !strings.HasSuffix(fn, Ext) {
		fn += Ext
	}
	if !strings.Contains(fn, "har/") {
		fn = path.Join("har", fn)
	}
	if !s.FS.Exists(fn) {
		return errors.Errorf("missing file [%s]", fn)
	}
	err := s.FS.Remove(fn, logger)
	if err != nil {
		return errors.Wrapf(err, "error deleting file [%s]", fn)
	}
	return nil
}

func (s *Service) Save(log *Log) error {
	fn := log.Key
	if !strings.HasSuffix(fn, Ext) {
		fn += Ext
	}
	if !strings.Contains(fn, "har/") {
		fn = path.Join("har", fn)
	}
	b := util.ToJSONBytes(&Wrapper{Log: log}, true)
	return s.FS.WriteFile(fn, b, filesystem.DefaultMode, true)
}

func (s *Service) Search(ctx context.Context, ps filter.ParamSet, q string, logger util.Logger) (result.Results, error) {
	return lo.FilterMap(s.List(logger), func(fn string, _ int) (*result.Result, bool) {
		log, err := s.Load(fn)
		if err != nil {
			logger.Warnf("error loading har [%s]: %+v", fn, err)
			return nil, false
		}
		res := result.NewResult("archive", log.Key, log.WebPath(), log.Key, "book", log, log, q)
		if len(res.Matches) > 0 {
			return res, true
		}
		return nil, false
	}), nil
}
