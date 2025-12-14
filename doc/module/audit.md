# Audit

The **`audit`** module provides a comprehensive audit framework for your application.
It enables detailed tracking of user actions, data changes, and system events with minimal performance overhead.

## Overview

This module provides enterprise-grade audit logging capabilities that integrate seamlessly with your application's existing data models and workflows.

**Key Features:**
- **Change Tracking**: Automatic diff generation for object modifications
- **Action Logging**: High-level user action tracking with metadata
- **Performance Optimized**: Asynchronous logging with minimal latency impact
- **Flexible Storage**: Works with any database supported by the `database` module
- **Web Interface**: Built-in admin interface for viewing and managing audit logs

## Core Components

### Audit Records
- **Audit**: High-level action tracking (app, action, client, server, user, metadata, timing)
- **Record**: Detailed change tracking (object type, primary key, field diffs, metadata)

### Service Layer
- **Transactional Support**: Atomic audit log creation with rollback capabilities
- **Automatic Diffing**: `ApplyObjSimple()` for hands-off object change tracking
- **Custom Records**: `Apply()` for manual audit record creation

### Web Interface
- **Admin Dashboard**: Complete CRUD operations via `/admin/audit` routes
- **Search & Filter**: Advanced filtering by user, action, timeframe, and object type
- **Detailed Views**: Drill-down from high-level audits to specific field changes

## Package Structure

### Controllers
- **`controller/audit/`** - Web interface for audit management
  - List, detail, create, edit, and delete audit records
  - Filtering and search capabilities
  - Export functionality for compliance reporting

### Libraries
- **`lib/audit/`** - Core audit service and data models
  - Audit and Record model definitions
  - Service layer for creating and querying audit logs
  - Integration helpers for automatic change tracking

## Dependencies

**Required Modules:**
- `database` - Provides database connectivity and migrations

**Recommended Modules:**
- `export` - Enables audit logs in generated code
- `user` - Provides user context for audit entries

## Configuration

The audit module automatically creates the necessary database tables during application startup. No additional configuration is required for basic functionality.

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/audit
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Database Module](database.md) - Required dependency for data storage
- [Export Module](export.md) - Recommended for audit log reporting
- [User Module](user.md) - Provides user context for audit entries
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
