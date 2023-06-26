package migrate

import (
	"time"

	"github.com/samber/lo"
)

type Migration struct {
	Idx     int       `json:"idx"`
	Title   string    `json:"title"`
	Src     string    `json:"src"`
	Created time.Time `json:"created"`
}

type migrationRow struct {
	Idx     int       `db:"idx"`
	Title   string    `db:"title"`
	Src     string    `db:"src"`
	Created time.Time `db:"created"`
}

func (r *migrationRow) toMigration() *Migration {
	return &Migration{
		Idx:     r.Idx,
		Title:   r.Title,
		Src:     r.Src,
		Created: r.Created,
	}
}

func toMigrations(rs []*migrationRow) Migrations {
	return lo.Map(rs, func(r *migrationRow, _ int) *Migration {
		return r.toMigration()
	})
}

type Migrations []*Migration

func (m Migrations) GetByIndex(idx int) *Migration {
	return lo.FindOrElse(m, nil, func(x *Migration) bool {
		return x.Idx == idx
	})
}
