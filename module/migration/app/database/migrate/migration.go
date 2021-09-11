package migrate

import "time"

type Migration struct {
	Idx     int       `json:"idx"`
	Title   string    `json:"title"`
	Src     string    `json:"src"`
	Created time.Time `json:"created"`
}

type Migrations []*Migration

type migrationDTO struct {
	Idx     int       `db:"idx"`
	Title   string    `db:"title"`
	Src     string    `db:"src"`
	Created time.Time `db:"created"`
}

func (dto *migrationDTO) toMigration() *Migration {
	return &Migration{
		Idx:     dto.Idx,
		Title:   dto.Title,
		Src:     dto.Src,
		Created: dto.Created,
	}
}

func toMigrations(dtos []migrationDTO) Migrations {
	ret := make(Migrations, 0, len(dtos))
	for _, dto := range dtos {
		ret = append(ret, dto.toMigration())
	}
	return ret
}
