package project

import (
	"github.com/samber/lo"
	"projectforge.dev/projectforge/app/lib/search/result"
)

func (s *Service) Search(q string) (result.Results, error) {
	ret := result.Results{}
	lo.ForEach(s.Projects(), func(p *Project, _ int) {
		if res := result.NewResult("project", p.Key, p.WebPath(), p.Title(), p.IconSafe(), p, p, q); len(res.Matches) > 0 {
			ret = append(ret, res)
		}
	})
	return ret, nil
}
