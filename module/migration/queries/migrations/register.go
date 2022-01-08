// $PF_IGNORE$
package migrations

import (
	"{{{ .Package }}}/app/lib/database/migrate"
)

func LoadMigrations(debug bool) {
	migrate.RegisterMigration("create initial database", Migration1InitialDatabase(debug))
}
