# [migration]

This is a module for [Project Forge](https://projectforge.dev). It provides database migrations and a common PostgreSQL or SQLite database

https://github.com/kyleu/projectforge/tree/master/module/migration

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

- Provides database migrations through a CLI command (`migrate`) and a web UI
- To run the migrations on normal startup, add `migrate.Migrate` to `InitApp`
- Migrations are defined in `./queries` SQL files, and registered by calling `AddMigration`
