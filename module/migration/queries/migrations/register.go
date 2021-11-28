// $PF_IGNORE$
package migrations

import (
	"{{{ .Package }}}/app/database/migrate"
)

func LoadMigrations() {
	migrate.RegisterMigration("create initial database", Migration1InitialDatabase())
}