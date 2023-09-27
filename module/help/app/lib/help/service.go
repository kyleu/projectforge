package help

import (
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/help"
)

type Service struct {
	Entries Entries
}

func NewService(logger util.Logger) (*Service, error) {
	entries, err := loadEntries(logger)
	if err != nil {
		return nil, errors.Wrap(err, "unable to load help entries")
	}
	logger.Debugf("loaded [%d] help entries", len(entries))
	return &Service{Entries: entries}, nil
}

func (s *Service) HTML(key string) string {
	ret := s.Entries.Get(key)
	if ret == nil {
		return ""
	}
	return ret.HTML
}

func (s *Service) Contains(key string) bool {
	return s.Entries.Get(key) != nil
}

func loadEntries(logger util.Logger) (Entries, error) {
	keys, err := help.List()
	if err != nil {
		logger.Errorf("unable to build documentation menu: %+v", err)
	}
	ret := make(Entries, 0, len(keys))
	for _, key := range keys {
		md, html, err := help.HTML(key)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load help file [%s]", key)
		}
		ret = append(ret, &Entry{Key: key, Markdown: md, HTML: html})
	}
	return ret, nil
}
