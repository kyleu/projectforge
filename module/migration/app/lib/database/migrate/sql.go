package migrate

import (
	"context"

	"github.com/pkg/errors"

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
	ret := make(MigrationFiles, 0, len(databaseMigrations))
	for _, x := range databaseMigrations {
		ret = append(ret, &MigrationFile{Title: x.Title, Content: x.Content})
	}
	return ret
}

func exec(ctx context.Context, file *MigrationFile, s *database.Service, logger util.Logger) (string, error) {
	sql := file.Content
	timer := util.TimerStart()
	logger.Infof("migration running SQL: %v", sql)
	_, err := s.Exec(ctx, sql, nil, -1, logger)
	if err != nil {
		return "", errors.Wrap(err, "cannot execute ["+file.Title+"]")
	}
	logger.Debugf("ran query [%s] in [%v]", file.Title, timer.EndString())
	return sql, nil
}
