package migrate

import "time"

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

func toMigrations(rs []migrationRow) Migrations {
	ret := make(Migrations, 0, len(rs))
	for _, r := range rs {
		ret = append(ret, r.toMigration())
	}
	return ret
}

type Migrations []*Migration

func (m Migrations) GetByIndex(idx int) *Migration {
	for _, x := range m {
		if x.Idx == idx {
			return x
		}
	}
	return nil
}
