# Database UI

The **`databaseui`** module provides a comprehensive web-based administration interface for [Project Forge](https://projectforge.dev) applications with database connectivity. It offers powerful database exploration, query execution, and monitoring capabilities through a modern, responsive web interface.

## Overview

This module is designed for applications that need **database administration capabilities** and provides:

- **Multi-Database Support**: Seamless integration with MySQL, PostgreSQL, SQLite, and SQL Server
- **Interactive SQL Editor**: Execute queries with syntax highlighting and real-time analysis
- **Database Exploration**: Browse tables, view schemas, and analyze data structures
- **Performance Monitoring**: Query tracing, execution timing, and performance diagnostics
- **Safe Operations**: Transaction controls with commit/rollback capabilities

**Security Notice**: This module is marked as **dangerous** and provides direct database access. Use only in trusted environments with proper authentication.

## Key Features

### Database Management
- Unified interface for multiple database connections
- Automatic database registration and discovery
- Connection health monitoring and diagnostics
- Support for multiple concurrent database sessions

### SQL Operations
- Query execution with real-time results
- Transaction management (commit/rollback)
- Query history and favorites
- Bulk operations and batch processing

### Performance & Monitoring
- Query execution timing and performance metrics
- Real-time tracing with configurable levels:
  - **None**: No tracing (production mode)
  - **Statements**: Log SQL statements only
  - **Results**: Include result set information
  - **Analysis**: Full query analysis and optimization hints
- Query profiling and optimization suggestions

### Data Exploration
- Table browsing with pagination and filtering
- Schema visualization and relationship mapping
- Row count statistics and table size analysis
- Data export capabilities (CSV, JSON, XML)

## Requirements

### Dependencies
- **`database`** module (required) - Provides core database connectivity
- Supported database drivers automatically included

### Supported Databases
- **MySQL** 5.7+ / MariaDB 10.3+
- **PostgreSQL** 11+
- **SQLite** 3.35+
- **SQL Server** 2017+

## Getting Started

### 1. Database Registration

Databases automatically register with the UI on startup. Ensure your database connections are configured in your application:

### 2. Access the Interface

Navigate to `/admin/database` in your application to access the database administration interface.

### 3. Basic Operations

- **Browse Tables**: Click on any database to explore its schema
- **Execute Queries**: Use the SQL editor to run custom queries
- **Monitor Performance**: Enable tracing to analyze query performance
- **Manage Transactions**: Use commit/rollback controls for safe operations

## Security Considerations

- **Authentication Required**: Always enable authentication for production use
- **Role-Based Access**: Implement proper role-based access controls
- **Query Validation**: Consider implementing query whitelisting for production
- **Connection Limits**: Monitor and limit concurrent database connections
- **Audit Logging**: Enable audit logging for all database operations

## Troubleshooting

### Common Issues

**Database Not Appearing**:
- Verify database module is properly configured
- Check database connection health
- Ensure proper permissions for schema access

**Query Execution Fails**:
- Check SQL syntax and permissions
- Verify transaction state
- Review error logs for connection issues

**Performance Issues**:
- Enable query tracing to identify bottlenecks
- Check connection pool settings
- Monitor database server resources

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/databaseui
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Database Module](database.md) - Core database connectivity
- [Configuration Variables](../running.md) - Available environment variables
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
