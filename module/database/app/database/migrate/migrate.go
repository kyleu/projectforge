package migrate

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/database"
)

func Migrate(ctx context.Context, s *database.Service, logger *zap.SugaredLogger) error {
	err := createMigrationTableIfNeeded(ctx, s, logger)
	if err != nil {
		return errors.Wrap(err, "unable to create migration table")
	}

	maxIdx := maxMigrationIdx(ctx, s, logger)

	if len(databaseMigrations) > maxIdx+1 {
		logger.Info(fmt.Sprintf("applying [%v] database migrations...", len(databaseMigrations)-maxIdx))
	}

	for i, file := range databaseMigrations {
		idx := i + 1
		switch {
		case idx == maxIdx:
			m := getMigrationByIdx(ctx, s, maxIdx, logger)
			if m == nil {
				continue
			}
			if m.Title != file.Title {
				logger.Info(fmt.Sprintf("migration [%v] name has changed from [%v] to [%v]", idx, m.Title, file.Title))
				err = removeMigrationByIdx(ctx, s, idx)
				if err != nil {
					return err
				}
				err = applyMigration(ctx, s, idx, file, logger)
				if err != nil {
					return err
				}
				continue
			}
			nc := file.Content
			if nc != m.Src {
				logger.Info(fmt.Sprintf("migration [%v:%v] content has changed from [%vB] to [%vB]", idx, file.Title, len(nc), len(m.Src)))
				err = removeMigrationByIdx(ctx, s, idx)
				if err != nil {
					return err
				}
				err = applyMigration(ctx, s, idx, file, logger)
				if err != nil {
					return err
				}
			}
		case idx > maxIdx:
			err = applyMigration(ctx, s, idx, file, logger)
			if err != nil {
				return err
			}
		default:
			// noop
		}
	}

	logger.Info(fmt.Sprintf("verified [%v] database migrations", maxIdx))

	return errors.Wrap(err, "error running database migration")
}

func applyMigration(ctx context.Context, s *database.Service, idx int, file *MigrationFile, logger *zap.SugaredLogger) error {
	logger.Info(fmt.Sprintf("applying database migration [%v]: %v", idx, file.Title))
	sql, err := exec(ctx, file, s, logger)
	if err != nil {
		return err
	}
	m := &Migration{Idx: idx, Title: file.Title, Src: sql, Created: time.Now()}
	return newMigration(ctx, s, m)
}
