// $PF_GENERATE_ONCE$
package migrations

import (
	"{{{ .Package }}}/app/lib/database/migrate"
)

func LoadMigrations(debug bool) {
	migrate.RegisterMigration("create initial database", Migration1InitialDatabase(debug))
}
