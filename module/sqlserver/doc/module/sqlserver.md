# SQL Server

The **`sqlserver`** module provides Microsoft SQL Server database integration for your application.
This module extends the base `database` module with SQL Server-specific functionality, enterprise features, and optimizations for Microsoft SQL Server environments.

## Overview

This module adds **Microsoft SQL Server database support** to your application and requires the `database` module. It provides:

- **Enterprise Database**: Full support for SQL Server's enterprise features
- **Driver Integration**: Official Microsoft SQL Server driver with connection pooling
- **Advanced Features**: Support for stored procedures, triggers, and SQL Server-specific syntax
- **Windows Integration**: Native Windows Authentication and Active Directory support

## Key Features

### SQL Server Compatibility
- Support for SQL Server 2017+ (including Azure SQL Database)
- Full T-SQL dialect support
- Advanced data types (geography, geometry, hierarchyid, etc.)
- Enterprise features (partitioning, compression, encryption)

### Enterprise Features
- **Stored Procedures**: Full support for complex stored procedures
- **User-Defined Types**: Custom data types and table-valued parameters
- **JSON Support**: Native JSON functions and operations
- **Spatial Data**: Geography and geometry data types

### Performance
- Optimized connection pooling for SQL Server workloads
- Query plan caching and prepared statement support
- Connection health monitoring with SQL Server-specific metrics
- Bulk operations and Table-Valued Parameters (TVP)

### Security
- Windows Authentication and Active Directory integration
- SSL/TLS encryption with certificate validation
- Row-level security and dynamic data masking
- Always Encrypted support for sensitive data

## Configuration

### Environment Variables

The module reads configuration from environment variables (with optional prefix):

- **`db_host`** - SQL Server hostname (default: `localhost`)
- **`db_port`** - SQL Server port (default: `1433`)
- **`db_user`** - Username for database connections
- **`db_password`** - Password for database connections (optional)
- **`db_database`** - Database name to connect to
- **`db_schema`** - Default schema to use (optional, default: `dbo`)
- **`db_max_connections`** - Maximum number of open connections
- **`db_debug`** - Enable SQL statement logging (`true`/`false`)

### Windows Authentication

```bash
# Use Windows Authentication (set db_user to local/native and omit password)
db_host=localhost\\SQLEXPRESS
db_database=MyDatabase
db_user=local
```

## Usage

### Basic Setup

```go
ctx := context.Background()
logger := util.RootLogger

// Open database connection from environment variables
svc, err := database.OpenDefaultSQLServer(ctx, logger)
if err != nil {
    return err
}
defer svc.Close()
```

Other examples below assume `ctx` and `logger` are already available.

### Custom Configuration

```go
// Manual configuration
params := &database.SQLServerParams{
    Host:     "sqlserver.example.com",
    Port:     1433,
    Username: "app_user",
    Password: "secure_password",
    Database: "app_database",
    Schema:   "app_schema",
    MaxConns: 25,
    Debug:    false,
}

svc, err := database.OpenSQLServerDatabase(ctx, "app", params, logger)
```

### Windows Authentication

```go
// Windows Authentication (no username/password)
params := &database.SQLServerParams{
    Host:     "localhost\\SQLEXPRESS",
    Database: "MyDatabase",
    Schema:   "dbo",
    Username: "local",
}

svc, err := database.OpenSQLServerDatabase(ctx, "MyDatabase", params, logger)
```

### With Environment Prefix

```go
// Use custom environment variable prefix
// Reads myapp_db_host, myapp_db_port, etc.
svc, err := database.OpenSQLServer(ctx, "mydb", "myapp_", logger)
```

## SQL Server-Specific Features

The examples below assume `svc`, `ctx`, and `logger` are already available.

### Stored Procedures

```go
// Execute stored procedure with parameters
_, err := svc.Exec(ctx, `exec GetUsersByRole @Role = ?`, nil, -1, logger, "admin")

// Stored procedure with output parameters
rows, err := svc.Query(ctx, `
    declare @Count int
    exec GetUserCount @Role = ?, @Count = @Count output
    select @Count`, nil, logger, "admin")
```

### Table-Valued Parameters

```go
// Create table type in SQL Server first:
// create type UserTableType as table (ID int, Name nvarchar(50))

// Use table-valued parameters for bulk operations
userTVP := mssql.TVP{
    TypeName: "UserTableType",
    Value: [][]interface{}{
        {1, "Alice"},
        {2, "Bob"},
    },
}

_, err := svc.Exec(ctx, `exec BulkInsertUsers @Users = ?`, nil, -1, logger, userTVP)
```

### JSON Operations

```go
// Insert JSON data
_, err := svc.Exec(ctx, `insert into users (data) values (?)`, nil, -1, logger,
    `{"name": "John", "preferences": {"theme": "dark"}}`)

// Query with JSON functions
rows, err := svc.Query(ctx,
    `select json_value(data, '$.name') from users where json_value(data, '$.active') = 'true'`,
    nil, logger)

// JSON path operations
rows, err := svc.Query(ctx,
    `select * from users where json_value(data, '$.preferences.theme') = 'dark'`,
    nil, logger)
```

### Spatial Data Types

```go
// Insert geometry data
_, err := svc.Exec(ctx,
    `insert into locations (point) values (geometry::Point(?, ?, 4326))`,
    nil, -1, logger, longitude, latitude)

// Spatial queries
rows, err := svc.Query(ctx, `
    select name from locations
    where point.STDistance(geometry::Point(?, ?, 4326)) < 1000`,
    nil, logger, userLong, userLat)
```

### MERGE Operations (Upsert)

```go
// SQL Server MERGE statement
_, err := svc.Exec(ctx, `
    merge users as target
    using (values (?, ?)) as source (email, name)
    on target.email = source.email
    when matched then
        update set name = source.name, updated_at = getdate()
    when not matched then
        insert (email, name, created_at) values (source.email, source.name, getdate());`,
    nil, -1, logger, email, name)
```

### Bulk Operations

```go
// Bulk copy for large data inserts
// Uses SqlBulkCopy equivalent functionality

txn, err := svc.StartTransaction(logger)
stmt, err := txn.Prepare(mssql.CopyIn("users", mssql.BulkOptions{}, "name", "email"))

for _, user := range largeUserList {
    _, err = stmt.Exec(user.Name, user.Email)
}

err = stmt.Close()
err = txn.Commit()
```

## Dependencies

This module requires:

- **`database`** - Core database functionality
- **[denisenkom/go-mssqldb](https://github.com/denisenkom/go-mssqldb)** - Microsoft SQL Server driver
- **[sqlx](https://github.com/jmoiron/sqlx)** - Enhanced SQL operations

## Production Considerations

### Connection Pooling
- Set `db_max_connections` based on SQL Server's `max connections` setting
- Monitor connection usage with built-in telemetry
- Consider using connection pooling middleware for high-traffic applications

### Performance Tuning
- Use SQL Server Management Studio for query optimization
- Enable query store for performance monitoring
- Configure appropriate memory and CPU settings
- Use columnstore indexes for analytical workloads

### High Availability
- Configure Always On Availability Groups
- Use read-only routing for read replicas
- Implement backup and disaster recovery strategies
- Monitor with SQL Server Agent jobs

### Licensing
- Understand SQL Server licensing requirements
- Consider Azure SQL Database for cloud deployments
- Plan for development, testing, and production environments

## Azure SQL Database

Special considerations for Azure SQL Database:

```go
// Azure SQL Database connection
params := &database.SQLServerParams{
    Host:     "myserver.database.windows.net",
    Port:     1433,
    Username: "myuser@myserver",
    Password: "mypassword",
    Database: "mydatabase",
}
```

### Azure-Specific Features
- Automatic scaling and performance tuning
- Built-in high availability
- Integrated security with Azure Active Directory
- Elastic pools for cost optimization

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/sqlserver
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Database Module](database.md) - Core database functionality
- [MySQL Module](mysql.md) - Popular open-source database
- [PostgreSQL Module](postgres.md) - Alternative enterprise database
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
