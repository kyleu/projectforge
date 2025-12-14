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
- WAL mode for improved concurrency
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

- **`db_file`** - SQLite database file path (relative to project root)
- **`db_schema`** - Default schema to use (optional)
- **`db_debug`** - Enable SQL statement logging (`true`/`false`)

### File Path Examples

```bash
# Relative to project root
db_file=data/app.db

# Absolute path
db_file=/var/lib/myapp/database.db

# Memory database (for testing)
db_file=:memory:

# Temporary database
db_file=""  # Uses temporary file
```

## Usage

### Basic Setup

```go
// Load configuration from environment
params := SQLiteParamsFromEnv("")

// Open database connection
db, err := database.OpenSQLiteDatabase(params)
if err != nil {
    return err
}
defer db.Close()
```

### Custom Configuration

```go
// Manual configuration
params := &SQLiteParams{
    File:   "data/myapp.db",
    Schema: "main",
    Debug:  false,
}

db, err := database.OpenSQLiteDatabase(params)
```

### Memory Database (Testing)

```go
// In-memory database for testing
params := &SQLiteParams{
    File:  ":memory:",
    Debug: true,
}

db, err := database.OpenSQLiteDatabase(params)
```

### With Environment Prefix

```go
// Use custom environment variable prefix
// Reads myapp_db_file, myapp_db_schema, etc.
params := SQLiteParamsFromEnv("myapp_")
```

## SQLite-Specific Features

### WAL Mode

SQLite databases automatically use WAL mode for better concurrency:

- WAL mode is enabled automatically
- Allows concurrent readers with single writer
- Better performance for web applications

### JSON Support

Modern SQLite includes JSON1 extension:

```go
// Insert JSON data
_, err := db.Exec(`insert into users (data) values (?)`, `{"name": "John", "preferences": {"theme": "dark"}}`)

// Query with JSON functions
rows, err := db.Query(`select json_extract(data, '$.name') from users`)

// JSON array operations
rows, err := db.Query(`select * from users where json_extract(data, '$.active') = 'true'`)
```

### Full-Text Search

SQLite's FTS5 extension provides powerful text search:

```go
// Create FTS table
_, err := db.Exec(`create virtual table posts_fts using fts5(title, content)`)

// Full-text search
rows, err := db.Query(`select * from posts_fts where posts_fts match 'database and sqlite'`)
```

### File Operations

```go
// Backup database
_, err := db.Exec(`vacuum into 'backup.db'`)

// Attach additional database
_, err := db.Exec(`attach database 'other.db' as other`)

// Database integrity check
rows, err := db.Query(`pragma integrity_check`)
```

## Development Workflow

### Database Creation

- Database file is created automatically on first connection
- No manual database creation required

### Testing

```go
// Use in-memory database for tests
func setupTestDB() *sql.DB {
    params := &SQLiteParams{File: ":memory:"}
    db, _ := database.OpenSQLiteDatabase(params)
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

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/sqlite
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Database Module](database.md) - Core database functionality
- [MySQL Module](mysql.md) - Popular database alternative
- [PostgreSQL Module](postgres.md) - Full-featured database option
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
