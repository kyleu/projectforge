# SQLite

This is a module for [Project Forge](https://projectforge.dev). It provides an API for accessing SQLite databases

https://github.com/kyleu/projectforge/tree/master/module/sqlite

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

- To configure a SQLite connection pool, call `database.OpenSQLiteDatabase`, passing `SQLiteParams`
- You can load the params from the environment by calling `SQLiteParamsFromEnv` (with optional prefix), this will read the following by default:
  - `DB_FILE` - SQLite database file to load, relative to project root
  - `DB_SCHEMA` - schema to use (optional)
  - `DB_DEBUG` - if set to `true`, will log all statements and parameters
