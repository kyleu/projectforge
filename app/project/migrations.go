package project

import (
	"context"
	"path/filepath"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/util"
)

type Migration struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type Migrations []*Migration

func (s *Service) MigrationList(ctx context.Context, prj *Project, logger util.Logger) (Migrations, error) {
	if !prj.HasModule("migration") {
		return nil, nil
	}
	fs, err := s.GetFilesystem(prj)
	if err != nil {
		return nil, err
	}
	const pth = "queries/migrations"
	if !fs.Exists(pth) {
		return nil, nil
	}
	files := fs.ListExtension(pth, "sql", nil, false, logger)
	return lo.Map(files, func(fn string, idx int) *Migration {
		content, _ := fs.ReadFile(filepath.Join(pth, fn))
		return &Migration{Filename: fn, Content: string(content)}
	}), nil
}
