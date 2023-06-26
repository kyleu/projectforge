package migrate

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/util"
)

type MigrationFile struct {
	Title   string
	Content string
}

type MigrationFiles []*MigrationFile

var databaseMigrations = MigrationFiles{}

func ClearMigrations() {
	databaseMigrations = MigrationFiles{}
}

func AddMigration(mf *MigrationFile) {
	databaseMigrations = append(databaseMigrations, mf)
}

func RegisterMigration(title string, content string) {
	AddMigration(&MigrationFile{Title: title, Content: content})
}

func GetMigrations() MigrationFiles {
	return lo.Map(databaseMigrations, func(x *MigrationFile, _ int) *MigrationFile {
		return &MigrationFile{Title: x.Title, Content: x.Content}
	})
}

func exec(ctx context.Context, file *MigrationFile, s *database.Service, tx *sqlx.Tx, logger util.Logger) (string, error) {
	sql := file.Content
	timer := util.TimerStart()
	logger.Infof("migration running SQL: %v", sql)
	_, err := s.Exec(ctx, sql, tx, -1, logger)
	if err != nil {
		return "", errors.Wrapf(err, "cannot execute [%s]", file.Title)
	}
	logger.Debugf("ran query [%s] in [%v]", file.Title, timer.EndString())
	return sql, nil
}
