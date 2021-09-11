package migrate

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/database"
	"{{{ .Package }}}/app/util"
)

func Migrate(ctx context.Context, s *database.Service, logger *zap.SugaredLogger) error {
	logger = logger.With("svc", "migrate")
	err := createMigrationTableIfNeeded(ctx, s, logger)
	if err != nil {
		return errors.Wrap(err, "unable to create migration table")
	}

	maxIdx := maxMigrationIdx(ctx, s, logger)

	if len(databaseMigrations) > maxIdx+1 {
		c := len(databaseMigrations) - maxIdx
		logger.Info(fmt.Sprintf("applying [%d] database %s...", c, util.PluralMaybe("migration", c)))
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
				logger.Info(fmt.Sprintf("migration [%d] name has changed from [%s] to [%s]", idx, m.Title, file.Title))
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
				logger.Info(fmt.Sprintf("migration [%d:%s] content has changed from [%dB] to [%dB]", idx, file.Title, len(nc), len(m.Src)))
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

	logger.Info(fmt.Sprintf("verified [%d] database %s", maxIdx, util.PluralMaybe("migration", maxIdx)))

	return errors.Wrap(err, "error running database migration")
}

func applyMigration(ctx context.Context, s *database.Service, idx int, file *MigrationFile, logger *zap.SugaredLogger) error {
	logger.Info(fmt.Sprintf("applying database migration [%d]: %s", idx, file.Title))
	sql, err := exec(ctx, file, s, logger)
	if err != nil {
		return err
	}
	m := &Migration{Idx: idx, Title: file.Title, Src: sql, Created: time.Now()}
	return newMigration(ctx, s, m)
}
