package migrate

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/database"
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

func exec(ctx context.Context, file *MigrationFile, s *database.Service, logger *zap.SugaredLogger) (string, error) {
	sql := file.Content
	startNanos := util.TimerStart()
	logger.Infof("migration running SQL: %v", sql)
	_, err := s.Exec(ctx, sql, nil, -1)
	if err != nil {
		return "", errors.Wrap(err, "cannot execute ["+file.Title+"]")
	}
	ms := util.MicrosToMillis(util.TimerEnd(startNanos))
	logger.Debug(fmt.Sprintf("ran query [%s] in [%v]", file.Title, ms))
	return sql, nil
}
