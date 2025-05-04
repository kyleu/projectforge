package help

import (
	"slices"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/help"
)

type Service struct {
	Entries Entries
}

func NewService(logger util.Logger) *Service {
	entries, err := loadEntries(logger)
	if err != nil {
		logger.Warnf("unable to load help entries: %+v", err)
	}
	logger.Debugf("loaded [%d] help entries", len(entries))
	return &Service{Entries: entries}
}

func (s *Service) Entry(key string) *Entry {
	return s.Entries.Get(key)
}

func (s *Service) Contains(key string) bool {
	return s.Entry(key) != nil
}

func loadEntries(logger util.Logger) (Entries, error) {
	keys, err := help.List()
	if err != nil {
		logger.Errorf("unable to build documentation menu: %+v", err)
	}
	var hd, ft string
	if slices.Contains(keys, "_header") {
		hd, _ = help.Content("_header.md")
	}
	if slices.Contains(keys, "_footer") {
		ft, _ = help.Content("_footer.md")
	}
	ret := make(Entries, 0, len(keys))
	for _, key := range keys {
		if key == "_header" || key == "_footer" {
			continue
		}
		ent, err := entryFor(key, hd, ft, 4)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to load help file [%s]", key)
		}
		if ent != nil && !strings.Contains(ent.Title, "!") {
			ret = append(ret, ent)
		}
	}
	return ret, nil
}

func entryFor(key string, hd string, ft string, indent int) (*Entry, error) {
	path := key
	if !strings.HasSuffix(path, util.ExtMarkdown) {
		path += util.ExtMarkdown
	}
	md, err := help.Content(path)
	if err != nil {
		return nil, err
	}
	var title string
	if strings.HasPrefix(md, "#") {
		if idx := strings.Index(md, "\n"); idx > -1 {
			title = strings.TrimSpace(strings.ReplaceAll(md[:idx], "#", ""))
			md = strings.TrimSpace(md[idx:])
		}
	}
	const dblLB = "\n\n"
	if hd != "" {
		md = hd + dblLB + md
	}
	if ft != "" {
		md += dblLB + ft
	}
	html := strings.TrimSpace(string(markdown.ToHTML([]byte(md), nil, nil)))
	html = util.StringJoin(util.StringSplitLinesIndented(html, indent, false, false), "\n")
	return &Entry{Key: key, Title: title, Markdown: md, HTML: html}, nil
}
