# PostgreSQL

The **`postgres`** module provides PostgreSQL database integration for your application.
This module extends the base `database` module with PostgreSQL-specific functionality, including advanced features, driver integration, and PostgreSQL-optimized performance.

## Overview

This module adds **PostgreSQL database support** to your application and requires the `database` module. It provides:

- **Advanced SQL Features**: Support for PostgreSQL's extended SQL capabilities
- **Driver Integration**: High-performance PostgreSQL driver with connection pooling
- **JSON/JSONB Support**: Native JSON document handling and querying
- **Configuration**: Environment-based configuration with PostgreSQL-specific optimizations

## Key Features

### PostgreSQL Compatibility
- Support for PostgreSQL 12+ (recommended 14+)
- Full PostgreSQL SQL dialect support
- Advanced data types (arrays, JSON, UUID, etc.)
- Window functions, CTEs, and advanced query features

### Advanced Data Types
- **JSON/JSONB**: Native document storage and querying
- **Arrays**: Native array data type support
- **UUID**: Built-in UUID generation and handling
- **Custom Types**: Support for enums and composite types

### Performance
- Optimized connection pooling for PostgreSQL workloads
- Query plan caching and prepared statement support
- Connection health monitoring with PostgreSQL-specific metrics
- Efficient bulk operations and COPY support

### Security
- SSL/TLS encryption with certificate validation
- Row-level security (RLS) support
- Advanced authentication methods
- SQL injection protection through prepared statements

## Configuration

### Environment Variables

The module reads configuration from environment variables (with optional prefix):

- **`db_host`** - PostgreSQL server hostname (default: `localhost`)
- **`db_port`** - PostgreSQL server port (default: `5432`)
- **`db_user`** - Username for database connections
- **`db_password`** - Password for database connections (optional)
- **`db_database`** - Database name to connect to
- **`db_schema`** - Default schema to use (optional, default: `public`)
- **`db_max_connections`** - Maximum number of active and idle connections
- **`db_debug`** - Enable SQL statement logging (`true`/`false`)

## Usage

### Basic Setup

```go
// Load configuration from environment
params := PostgresParamsFromEnv("")

// Open database connection
db, err := database.OpenPostgresDatabase(params)
if err != nil {
    return err
}
defer db.Close()
```

### Custom Configuration

```go
// Manual configuration
params := &PostgresParams{
    Host:           "postgres.example.com",
    Port:           5432,
    User:           "app_user",
    Password:       "secure_password",
    Database:       "app_database",
    Schema:         "app_schema",
    MaxConnections: 25,
    Debug:          false,
}

db, err := database.OpenPostgresDatabase(params)
```

### With Environment Prefix

```go
// Use custom environment variable prefix
// Reads myapp_db_host, myapp_db_port, etc.
params := PostgresParamsFromEnv("myapp_")
```

## Dependencies

This module requires:

- **`database`** - Core database functionality
- **[lib/pq](https://github.com/lib/pq)** - Pure Go PostgreSQL driver
- **[sqlx](https://github.com/jmoiron/sqlx)** - Enhanced SQL operations

## Production Considerations

### Connection Pooling
- Set `db_max_connections` based on PostgreSQL's `max_connections` setting
- Monitor connection usage with built-in telemetry
- Consider connection poolers like PgBouncer for high-concurrency applications

### Performance Tuning
- Use `explain analyze` for query optimization
- Configure appropriate work_mem and shared_buffers
- Enable query statistics collection
- Use partial indexes for large tables

### Schema Management
- Use the [migration](migrations.md) module for automatic schema changes
- Leverage PostgreSQL's transactional DDL
- Consider schema-per-tenant for multi-tenant applications

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/postgres
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Database Module](database.md) - Core database functionality
- [MySQL Module](mysql.md) - MySQL database support
- [SQLite Module](sqlite.md) - Embedded database option
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
