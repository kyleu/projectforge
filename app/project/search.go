package project

import (
	"context"

	"go.uber.org/zap"
	"projectforge.dev/projectforge/app/lib/search/result"
)

func (s *Service) Search(ctx context.Context, q string, logger *zap.SugaredLogger) (result.Results, error) {
	ret := result.Results{}
	for _, p := range s.Projects() {
		if res := result.NewResult("project", p.Key, p.WebPath(), p.Title(), p.IconSafe(), p, q); len(res.Matches) > 0 {
			ret = append(ret, res)
		}
	}
	return ret, nil
}
