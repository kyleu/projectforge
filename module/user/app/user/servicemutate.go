// $PF_GENERATE_ONCE$
package user

import (
	"context"
	"time"

	"github.com/google/uuid"

	"{{{ .Package }}}/app/lib/filesystem"
	"{{{ .Package }}}/app/util"
)

func (s *Service) Create(ctx context.Context, logger util.Logger, _ any, models ...*User) error {
	return s.Save(ctx, logger, models...)
}

func (s *Service) CreateIfNeeded(ctx context.Context, userID uuid.UUID, username string, _ any, logger util.Logger) error {
	if curr, _ := s.Get(ctx, nil, userID, logger); curr == nil {
		err := s.Create(ctx, logger, &User{ID: userID, Name: username, Created: time.Now()})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Update(ctx context.Context, _ any, model *User, logger util.Logger) error {
	return s.Save(ctx, logger, model)
}

func (s *Service) Save(ctx context.Context, logger util.Logger, models ...*User) error {
	for _, model := range models {
		b := util.ToJSONBytes(model, true)
		err := s.files.WriteFile(dirFor(model.ID), b, filesystem.DefaultMode, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, logger util.Logger) error {
	return s.files.Remove(dirFor(id), logger)
}
