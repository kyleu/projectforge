package export

import (
	"context"
	"github.com/pkg/errors"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/tools/pfdb/app/database"
)

type Service struct {
	Logger util.Logger
	DB     *database.Service
}

func NewService(logger util.Logger) (*Service, error) {
	return &Service{Logger: logger}, nil
}

func (s *Service) Run() error {
	s.Logger.Info("starting [pfdb]...")
	db, err := s.loadDatabase()
	if err != nil {
		return err
	}
	s.DB = db
	s.Logger.Infof("successfully connected to [%s]", s.DB.Type.Title)
	return nil
}

func (s *Service) loadDatabase() (*database.Service, error) {
	key := util.GetEnv("key", util.AppKey)
	switch dbType := util.GetEnv("database", "postgres"); dbType {
	case "mysql":
		return nil, nil
	case "postgres":
		return database.OpenPostgres(context.Background(), key, "", s.Logger)
	case "sqlite":
		return nil, nil
	case "sqlserver":
		return nil, nil
	default:
		return nil, errors.Errorf("invalid database type [%s]", dbType)
	}
}
