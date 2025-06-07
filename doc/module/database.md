# Database

The **`database`** module provides foundational database access and management capabilities for [Project Forge](https://projectforge.dev) applications. It serves as the core abstraction layer for relational database operations with support for multiple database engines.

## Overview

This module provides **engine-agnostic database functionality** and requires a database engine module (such as `postgres`, `mysql`, `sqlite`, or `sqlserver`) to function. It offers:

- **Service Infrastructure**: Database connection management with health monitoring
- **Query Framework**: Type-safe SQL execution with quicktemplate integration
- **Transaction Support**: Comprehensive transaction handling and connection pooling
- **Observability**: Built-in telemetry, metrics, and performance monitoring

## Key Features

### Engine Flexibility
- Support for PostgreSQL, MySQL, SQLite, and SQL Server
- Engine-agnostic query interface
- Seamless engine switching for development and production

### Performance
- Connection pooling via [sqlx](https://github.com/jmoiron/sqlx)
- Query compilation and caching
- Zero-allocation database utilities where possible
- Built-in performance metrics and profiling

### Developer Experience
- Type-safe database operations
- Quicktemplate SQL query compilation
- Comprehensive error handling with context
- Debug logging and query tracing

### Observability
- OpenTelemetry database span tracking
- Connection pool and query performance metrics
- Health check integration
- Detailed error reporting and logging

## Package Structure

### Core Infrastructure

- **`service/`** - Database service management and lifecycle
  - Connection pool management
  - Health check implementations
  - Service registration and discovery
  - Configuration management

### Libraries

- **`lib/database/`** - Core database utilities and abstractions
  - **Connection Management**: Pool configuration and lifecycle
  - **Query Execution**: Type-safe SQL operations (Select, Get, Query, Exec)
  - **Transaction Handling**: Context-aware transaction management
  - **Telemetry Integration**: Spans, metrics, and performance tracking
  - **Error Handling**: Detailed error context and recovery

### Query Framework

- **`queries/`** - SQL query templates (compiled with quicktemplate)
  - Parameterized query support
  - Template inheritance and composition
  - Type-safe parameter binding
  - Query result mapping

## Database Operations

### Basic Queries

```go
// Select multiple rows
users, err := db.Select(ctx, "select * from users where active = ?", true)

// Get single row
user, err := db.Get(ctx, "select * from users where id = ?", userID)

// Execute statement
result, err := db.Exec(ctx, "update users set active = ? where id = ?", false, userID)
```

### SQL Helpers

The database module provides SQL query building helpers for common operations:

#### Select Queries
- `SQLSelect(columns, tables, where, orderBy, limit, offset, dbt)` - Basic SELECT query
- `SQLSelectSimple(columns, tables, dbt, where...)` - Simple SELECT with optional WHERE conditions
- `SQLSelectGrouped(columns, tables, where, groupBy, orderBy, limit, offset, dbt)` - SELECT with GROUP BY support

#### Insert Operations
- `SQLInsert(table, columns, rows, dbt)` - INSERT statement for one or more rows
- `SQLInsertReturning(table, columns, rows, returning, dbt)` - INSERT with RETURNING clause

#### Update Operations
- `SQLUpdate(table, columns, where, dbt)` - UPDATE statement with WHERE clause
- `SQLUpdateReturning(table, columns, where, returned, dbt)` - UPDATE with RETURNING clause

#### Upsert Operations
- `SQLUpsert(table, columns, rows, conflicts, updates, dbt)` - INSERT with ON CONFLICT (PostgreSQL) or MERGE (SQL Server)

#### Delete Operations
- `SQLDelete(table, where, dbt)` - DELETE statement with mandatory WHERE clause

#### Utility Functions
- `SQLInClause(column, numParams, offset, dbt)` - Generate IN clause with placeholders

All functions handle database-specific placeholder syntax and SQL dialect differences automatically.

### Transaction Management

```go
// Automatic transaction handling
err := db.WithTransaction(ctx, func(tx *sql.Tx) error {
    // Multiple operations in transaction
    _, err := tx.Exec("INSERT INTO ...")
    return err
})
```

### Query Templates

SQL files in `queries/` directory are compiled with quicktemplate:

```sql
-- queries/user.sql
SELECT {%s= columns %} FROM users 
WHERE active = {%v active %}
{% if search %}AND name LIKE {%v search %}{% endif %}
ORDER BY created_at DESC
```

## Configuration

The database module supports configuration through environment variables and database engine modules:

### Connection Settings
- Database-specific connection strings (engine-dependent)
- Connection pool size and timeout settings
- SSL/TLS configuration for secure connections

### Performance Tuning
- `db_pool_max_open` - Maximum open connections
- `db_pool_max_idle` - Maximum idle connections
- `db_pool_max_lifetime` - Connection maximum lifetime

### Observability
- `db_debug_enabled` - Enable query logging and debugging
- Database metrics collection (automatic via telemetry)
- Connection health monitoring

## Required Engine Module

The `database` module requires one of the following engine-specific modules:

- **`postgres`** - PostgreSQL support with advanced features
- **`mysql`** - MySQL/MariaDB compatibility
- **`sqlite`** - Embedded SQLite database
- **`sqlserver`** - Microsoft SQL Server integration

Each engine module provides:
- Driver-specific connection handling
- Database schema migration support
- Engine-optimized query patterns
- Platform-specific configuration

## Dependencies

The `database` module integrates with:

- **[sqlx](https://github.com/jmoiron/sqlx)** - Enhanced SQL database driver
- **[OpenTelemetry](https://opentelemetry.io/)** - Database operation tracing
- **[Quicktemplate](https://github.com/valyala/quicktemplate)** - SQL query compilation

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/database
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [MySQL Module](mysql.md) - MySQL/MariaDB integration
- [PostgreSQL Module](postgres.md) - PostgreSQL-specific functionality
- [SQLite Module](sqlite.md) - Embedded database support
- [SQL Server Module](sqlserver.md) - Microsoft SQL Server integration
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
