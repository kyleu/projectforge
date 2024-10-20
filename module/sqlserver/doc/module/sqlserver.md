# SQL Server

This is a module for [Project Forge](https://projectforge.dev). It provides an API for accessing Microsoft SQL Server databases.

https://github.com/kyleu/projectforge/tree/main/module/sqlserver

### License

Licensed under [CC0](https://creativecommons.org/publicdomain/zero/1.0)

### Usage

- To configure a SQL Server connection pool, call `database.OpenSQLServerDatabase`, passing `SQLServerParams`
- You can load the params from the environment by calling `SQLServerParamsFromEnv` (with optional prefix), this will read the following by default:
  - `DB_HOST` - hostname to use, defaults to `localhost`
  - `DB_PORT` - port to use, defaults to 1433
  - `DB_USER` - username for connections
  - `DB_PASSWORD` - password for connections (optional)
  - `DB_DATABASE` - database to use
  - `DB_SCHEMA` - schema to use (optional)
  - `DB_MAX_CONNECTIONS` - max active and idle connections
  - `DB_DEBUG` - if set to `true`, will log all statements and parameters
