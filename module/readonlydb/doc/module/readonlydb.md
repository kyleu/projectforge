# Read-only Database

The **`readonlydb`** module provides a separate read-only database connection for your application.
This enables read-heavy operations to be offloaded to a dedicated read replica or separate database instance, improving performance and reducing load on the primary database.

## Overview

This module extends applications with database module enabled by adding:

- **Separate Read Connection**: Independent database connection pool for read operations
- **Flexible Configuration**: Complete control over read database settings
- **Performance Optimization**: Offload read queries to reduce primary database load
- **High Availability**: Support for read replicas and database separation

## Key Features

### Performance Benefits
- Dedicated connection pool for read operations
- Reduces load on primary write database
- Optimizes query distribution across database infrastructure
- Configurable connection limits for read workloads

### Configuration Flexibility
- Independent database server configuration
- Separate credentials and connection settings
- Support for different database schemas
- Debug logging for read operations

### Database Support
- **PostgreSQL**: Full support for read replicas
- **MySQL**: Compatible with MySQL read slaves
- **SQLite**: Separate database file support
- **SQL Server**: Read-only connection support

## Configuration Variables

The module provides comprehensive configuration through environment variables:

### Connection Settings
- **`read_db_host`** (string): Hostname for read database connection (default: `localhost`)
- **`read_db_port`** (int): Port for read database connection (default: `3306`)
- **`read_db_user`** (string): Username for read database authentication
- **`read_db_password`** (string): Password for read database authentication
- **`read_db_database`** (string): Database name for read connections
- **`read_db_schema`** (string): Schema name for read operations
- **`read_db_debug`** (bool): Enable debug logging for all read database interactions

## Usage

### Application Integration

When enabled, the read-only database connection is available through the application state:

```go
// Access read-only database connection
readDB := as.Services.ReadDB

// Use for read operations
rows, err := readDB.Query("SELECT * FROM users WHERE active = true")
```

### Best Practices

- Use read connection for reports and analytics
- Route heavy SELECT queries to read database
- Keep write operations on primary database
- Monitor read replica lag in production

## Dependencies

- **Required**: [`database`](database.md) - Provides core database functionality
- **Compatible**: Works with all supported database technologies

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/readonlydb
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Database Module](database.md) - Core database functionality
- [Configuration Variables](../running.md) - Complete environment variable reference
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
