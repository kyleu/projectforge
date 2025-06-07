# Database Migration

The **`migration`** module provides comprehensive database migration capabilities for [Project Forge](https://projectforge.dev) applications. It enables structured database schema evolution with version control and automated migration execution.

## Overview

This module provides:

- **CLI Migration Command**: Execute migrations from the command line using `./bin/app migrate`
- **Web UI Interface**: View and manage migrations through the web interface
- **Automatic Migration**: Run migrations on application startup
- **Version Control**: Track schema changes with versioned migration files
- **Multi-Database Support**: Works with PostgreSQL, SQLite, MySQL, and SQL Server

## Key Features

### Migration Management
- Versioned migration files in SQL format
- Automatic migration tracking and execution
- Rollback capabilities for schema changes
- Migration status monitoring and reporting

### Integration Options
- **Command Line**: `./bin/app migrate` for manual execution
- **Web Interface**: Browser-based migration management
- **Application Startup**: Automatic migration on app initialization
- **CI/CD Integration**: Execute migrations in deployment pipelines

### Database Support
- PostgreSQL (primary support)
- SQLite (embedded applications)
- MySQL (enterprise applications)
- SQL Server (enterprise applications)

## Package Structure

### Migration Library

- **`lib/database/migrate/`** - Core migration functionality
  - Migration file parsing and execution
  - Version tracking and state management
  - Database schema validation
  - Rollback and recovery operations

### Migration Files

- **`queries/migrations/`** - SQL migration files directory
  - Strongly-typed `quicktemplate` functions
  - Versioned schema changes
  - Data transformation scripts
  - Index and constraint definitions
  - Seed data initialization

## Usage

### Basic Setup

1. **Add Migration Dependency**: Include `migration` in your project modules
2. **Create Migration Files**: Add SQL files to `./queries/migrations/`
3. **Register Migrations**: Call `AddMigration()` to register migration files
4. **Execute Migrations**: Use CLI command or automatic startup

### Creating Migrations

Create SQL files in `./queries/migrations/` directory:

```sql
-- 001_create_users_table.sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

### Registration

Register migrations in your application:

```go
// In your application initialization
func init() {
    migrations.AddMigration("001_create_users_table.sql", createUsersTable)
    migrations.AddMigration("002_add_user_roles.sql", addUserRoles)
}
```

### Execution Methods

#### Command Line Execution
```bash
# Execute all pending migrations
./bin/app migrate
```

#### Automatic Startup Execution
```go
// Add to your InitApp function
func InitApp(ctx context.Context, st *app.State) error {
    // Run migrations on startup
    err := migrate.Migrate(ctx, st.DB, st.Logger)
    if err != nil {
        return errors.Wrap(err, "migration failed")
    }
    return nil
}
```

#### Web Interface
- Navigate to `/admin/migrations` in your application
- View migration status and execution history

## Dependencies

- **`database`** module - Database connection and management

## Best Practices

### Migration File Structure
- One logical change per migration file
- Test migrations on sample data
- Use transactions for atomic operations

### Version Control
- Commit migration files with application code
- Never modify existing migration files
- Create new migrations for schema changes
- Tag releases after successful migrations

### Production Deployment
- Test migrations in staging environments
- Backup databases before migration execution
- Monitor migration execution time and performance
- Plan for rollback scenarios

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/migration
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Database Module](database.md) - Database connection and management
- [Configuration Guide](../running.md) - Environment variable configuration
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
