# SQLite

The **`sqlite`** module provides SQLite database integration for your application.
This module extends the base `database` module with SQLite-specific functionality, offering an embedded database solution perfect for development, testing, and lightweight production deployments.

## Overview

This module adds **SQLite embedded database support** to your application and requires the `database` module. It provides:

- **Embedded Database**: Zero-configuration embedded database with no external dependencies
- **File-Based Storage**: Single-file database storage with backup and portability
- **Development-Friendly**: Perfect for development, testing, and prototyping
- **Production-Ready**: Suitable for lightweight production applications

## Key Features

### Embedded Database
- No external database server required
- Single-file database storage
- Zero-configuration setup
- Cross-platform compatibility

### Performance
- In-process database access (no network overhead)
- `busy_timeout` set to 10 seconds to reduce lock errors
- Connection pooling for multi-threaded access
- Optimized for read-heavy workloads

### Portability
- Single file contains entire database
- Easy backup and migration
- Version control friendly (for small databases)
- Platform-independent database files

### Development Benefits
- Instant setup for new projects
- No database server installation required
- Perfect for unit testing and CI/CD
- Simplified deployment process

## Configuration

### Environment Variables

The module reads configuration from environment variables (with optional prefix):

- **`db_file`** - SQLite database file path (defaults to `<config_dir>/<app_key>.sqlite`)
- **`db_user`** - Optional auth username (requires `db_password`)
- **`db_password`** - Optional auth password (requires `db_user`)
- **`db_schema`** - Default schema to use (optional, defaults to `public`, SQLite typically uses `main`)
- **`db_debug`** - Enable SQL statement logging (`true`/`false`)

When using `SQLiteParamsFromEnv`, the default file path is derived from `<config_dir>/<app_key>` and `.sqlite` is appended if missing.

### File Path Examples

```bash
# Default (if db_file is unset)
# <config_dir>/<app_key>.sqlite

# Relative path (resolved from the current working directory)
db_file=data/app.sqlite

# Absolute path
db_file=/var/lib/myapp/app.sqlite

# Memory database (for testing)
db_file=:memory:
```

## Usage

### Basic Setup

```go
logger := util.RootLogger
ctx := context.Background()

// Open database connection using environment defaults
db, err := database.OpenDefaultSQLite(ctx, logger)
if err != nil {
    return err
}
defer db.Close()
```

### Custom Configuration

```go
// Manual configuration
params := &SQLiteParams{
    File:   "data/myapp.sqlite",
    Schema: "main",
    Debug:  false,
}

db, err := database.OpenSQLiteDatabase(ctx, "myapp", params, logger)
```

### Memory Database (Testing)

```go
// In-memory database for testing
params := &SQLiteParams{
    File:  ":memory:",
    Debug: true,
}

db, err := database.OpenSQLiteDatabase(ctx, "test", params, logger)
```

### With Environment Prefix

```go
// Use custom environment variable prefix
// Reads myapp_db_file, myapp_db_schema, etc.
params := SQLiteParamsFromEnv("myapp", "myapp_")
db, err := database.OpenSQLiteDatabase(ctx, "myapp", params, logger)
```

## SQLite-Specific Features

### Connection Pragmas

SQLite connections are opened with:

- `foreign_keys=1`
- `busy_timeout=10000` (10 seconds)
- `trusted_schema=0`

To enable WAL mode, run `pragma journal_mode = wal` once per database.

### JSON Support

Modern SQLite includes JSON1 extension:

```go
// Insert JSON data
_, err := db.Exec(ctx, `insert into users (data) values (?)`, nil, 1, logger, `{"name": "John", "preferences": {"theme": "dark"}}`)

// Query with JSON functions
rows, err := db.Query(ctx, `select json_extract(data, '$.name') from users`, nil, logger)

// JSON array operations
rows, err := db.Query(ctx, `select * from users where json_extract(data, '$.active') = 'true'`, nil, logger)
```

### Full-Text Search

SQLite's FTS5 extension provides powerful text search:

```go
// Create FTS table
_, err := db.Exec(ctx, `create virtual table posts_fts using fts5(title, content)`, nil, -1, logger)

// Full-text search
rows, err := db.Query(ctx, `select * from posts_fts where posts_fts match 'database and sqlite'`, nil, logger)
```

### File Operations

```go
// Backup database
_, err := db.Exec(ctx, `vacuum into 'backup.db'`, nil, -1, logger)

// Attach additional database
_, err := db.Exec(ctx, `attach database 'other.db' as other`, nil, -1, logger)

// Database integrity check
rows, err := db.Query(ctx, `pragma integrity_check`, nil, logger)
```

## Development Workflow

### Database Creation

- Database file is created automatically on first connection
- No manual database creation required

### Testing

```go
// Use in-memory database for tests
func setupTestDB() *database.Service {
    logger := util.RootLogger
    ctx := context.Background()
    params := &SQLiteParams{File: ":memory:"}
    db, _ := database.OpenSQLiteDatabase(ctx, "test", params, logger)
    return db
}
```

## Dependencies

This module requires:

- **`database`** - Core database functionality
- **[modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)** - Pure Go SQLite driver
- **[sqlx](https://github.com/jmoiron/sqlx)** - Enhanced SQL operations

## Production Considerations

### File Management
- Ensure adequate disk space for database growth
- Implement regular backup procedures
- Consider file permissions and security
- Monitor database file size

### Performance Optimization
- Use appropriate indexes for query patterns
- Consider VACUUM operations for maintenance
- Enable query planner statistics
- Use prepared statements for repeated queries

### Limitations
- Single writer, multiple readers concurrency model
- Database size practical limit (~1TB)
- No built-in replication or clustering
- Limited concurrent write performance

## Build Support

The SQLite module is enabled for darwin, linux (386/amd64/arm/arm64/riscv64), and windows/amd64 builds.
Other targets use a stub that returns "SQLite is not enabled in this build".

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/sqlite
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Database Module](database.md) - Core database functionality
- [MySQL Module](mysql.md) - Popular database alternative
- [PostgreSQL Module](postgres.md) - Full-featured database option
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
