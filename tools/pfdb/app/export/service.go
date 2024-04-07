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

func (s *Service) Run(dbType string, key string) error {
	s.Logger.Info("starting [pfdb]...")
	db, err := s.loadDatabase(dbType, key)
	if err != nil {
		return err
	}
	s.DB = db
	s.Logger.Infof("successfully connected to [%s]", s.DB.Type.Title)
	return nil
}

func (s *Service) loadDatabase(dbType string, dbKey string) (*database.Service, error) {
	key := util.GetEnv("database_key", dbKey)
	dbType = util.GetEnv("database_type", dbType)
	switch dbType {
	case "mysql":
		return database.OpenMySQL(context.Background(), key, "", s.Logger)
	case "postgres":
		return database.OpenPostgres(context.Background(), key, "", s.Logger)
	case "sqlite":
		return database.OpenSQLite(context.Background(), key, "", s.Logger)
	case "sqlserver":
		return database.OpenSQLServer(context.Background(), key, "", s.Logger)
	default:
		return nil, errors.Errorf("invalid database type [%s]", dbType)
	}
}
