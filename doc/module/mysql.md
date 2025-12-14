# MySQL

The **`mysql`** module provides MySQL database integration for your application.
This module extends the base `database` module with MySQL-specific functionality, including driver integration, connection management, and MySQL-optimized features.

## Overview

This module adds **MySQL/MariaDB database support** to your application and requires the `database` module. It provides:

- **Driver Integration**: Official MySQL driver with connection pooling
- **MySQL Features**: Support for MySQL-specific SQL syntax and features
- **Performance Optimization**: Connection pooling and query optimization for MySQL
- **Configuration**: Environment-based configuration with sensible defaults

## Key Features

### MySQL Compatibility
- Support for MySQL 5.7+ and MariaDB 10.3+
- MySQL-specific SQL dialect handling
- Full UTF8MB4 character set support
- Native JSON data type support

### Performance
- Optimized connection pooling for MySQL workloads
- Query plan caching and prepared statement support
- Connection health monitoring
- Configurable timeouts and limits

### Security
- SSL/TLS encryption support
- Secure password handling
- Connection validation and health checks
- SQL injection protection through prepared statements

## Configuration

### Environment Variables

The module reads configuration from environment variables (with optional prefix):

- **`db_host`** - MySQL server hostname (default: `localhost`)
- **`db_port`** - MySQL server port (default: `3306`)
- **`db_user`** - Username for database connections
- **`db_password`** - Password for database connections (optional)
- **`db_database`** - Database name to connect to
- **`db_schema`** - Default schema to use (optional)
- **`db_max_connections`** - Maximum number of active and idle connections
- **`db_debug`** - Enable SQL statement logging (`true`/`false`)

## Usage

### Basic Setup

```go
// Load configuration from environment
params := MySQLParamsFromEnv("")

// Open database connection
db, err := database.OpenMySQLDatabase(params)
if err != nil {
    return err
}
defer db.Close()
```

### Custom Configuration

```go
// Manual configuration
params := &MySQLParams{
    Host:           "mysql.example.com",
    Port:           3306,
    User:           "app_user",
    Password:       "secure_password",
    Database:       "app_database",
    MaxConnections: 25,
    Debug:          false,
}

db, err := database.OpenMySQLDatabase(params)
```

### With Environment Prefix

```go
// Use custom environment variable prefix
// Reads myapp_db_host, myapp_db_port, etc.
params := MySQLParamsFromEnv("myapp_")
```

## MySQL-Specific Features

### JSON Support

MySQL's native JSON data type is fully supported:

```go
// Insert JSON data
_, err := db.Exec(`insert into users (data) values (?)`,
    `{"name": "John", "preferences": {"theme": "dark"}}`)

// Query with JSON functions
rows, err := db.Query(`select json_extract(data, '$.name') from users`)
```

### Character Set Handling

The module automatically configures UTF8MB4 for full Unicode support:

```go
// Automatic UTF8MB4 configuration
// Supports emojis and 4-byte Unicode characters
```

## Dependencies

This module requires:

- **`database`** - Core database functionality
- **[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)** - Official MySQL driver
- **[sqlx](https://github.com/jmoiron/sqlx)** - Enhanced SQL operations

## Production Considerations

### Connection Pooling
- Set `db_max_connections` based on MySQL's `max_connections` setting
- Monitor connection usage with built-in telemetry
- Consider read/write splitting for high-traffic applications

### Performance Tuning
- Enable query caching in MySQL configuration
- Use prepared statements for repeated queries
- Monitor slow query logs

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/mysql
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Database Module](database.md) - Core database functionality
- [PostgreSQL Module](postgres.md) - PostgreSQL database support
- [SQLite Module](sqlite.md) - Embedded database option
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
