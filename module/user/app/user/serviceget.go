// $PF_GENERATE_ONCE$
package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/util"
)

func (s *Service) List(ctx context.Context, params *filter.Params, logger util.Logger) (Users, error) {
	files := s.files.ListJSON("users", nil, true, s.logger)
	ret := make(Users, 0, len(files))
	for idx, idStr := range files {
		id := util.UUIDFromStringOK(idStr)
		r, err := s.Get(ctx, nil, id, logger)
		if err != nil {
			return nil, errors.Wrap(err, "error loading request ["+idStr+"]")
		}
		ret = append(ret, r)
		if (params == nil && idx > 100) || idx >= params.Limit {
			break
		}
	}
	return ret, nil

}

func (s *Service) Count(ctx context.Context, whereClause string, logger util.Logger, args ...any) (int, error) {
	files := s.files.ListJSON("users", nil, true, s.logger)
	return len(files), nil
}

func (s *Service) Get(ctx context.Context, _ any, id uuid.UUID, logger util.Logger) (*User, error) {
	content, err := s.files.ReadFile(dirFor(id))
	if err != nil {
		return nil, err
	}
	var ret *User
	err = util.FromJSON(content, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *Service) GetMultiple(ctx context.Context, logger util.Logger, ids ...uuid.UUID) (Users, error) {
	return lo.Map(ids, func(id uuid.UUID, _ int) *User {
		ret, _ := s.Get(ctx, nil, id, logger)
		return ret
	}), nil
}

func (s *Service) Search(ctx context.Context, query string, params *filter.Params, logger util.Logger) (Users, error) {
	return nil, nil
}
