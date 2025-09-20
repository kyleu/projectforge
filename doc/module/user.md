# User

The **`user`** module provides user management functionality for [Project Forge](https://projectforge.dev) applications. It handles user authentication, profiles, session management, and provides both filesystem-based and database-backed user storage options.

## Overview

This module provides:

- **User Profile Management**: User account creation, updates, and profile handling
- **Session Integration**: Seamless integration with Project Forge's session system
- **Flexible Storage**: Support for both filesystem-based and database-backed user storage
- **Authentication Support**: Works with OAuth and other authentication modules
- **Customizable Models**: Easily customizable user models through export definitions

## Key Features

### User Management
- User profile creation and management
- Automatic user session integration
- Profile persistence across sessions
- User data validation and sanitization

### Storage Options
- **Filesystem**: Default implementation using filesystem storage
- **Database**: Full database integration with customizable schema
- **Hybrid**: Combine both approaches as needed

### Integration
- Works seamlessly with `oauth` module for authentication
- Integrates with `database` module for persistent storage
- Compatible with `audit` module for user activity tracking
- Supports role-based access control when combined with RBAC modules

## Package Structure

### Core Components

- **`user/`** - User model definitions and core logic
  - User struct definitions
  - Profile management utilities
  - Session integration handlers

- **`controller/`** - HTTP handlers for user operations
  - Profile viewing and editing
  - User registration and onboarding
  - Session management endpoints

- **`service/`** - User business logic and operations
  - User creation and updates
  - Profile validation
  - User lookup and search functionality

### Storage Implementations

- **Filesystem Storage**: Default implementation for simple deployments
  - JSON-based user profiles
  - File-based session persistence
  - No external dependencies

- **Database Storage**: Full database integration (requires `database` module)
  - Customizable user schema
  - SQL-based queries and operations
  - Transaction support

## Configuration

### Default User Model

The module provides a basic user model with required fields:

```go
type User struct {
    ID      uuid.UUID  `json:"id"`
    Name    string     `json:"name"`
    Created time.Time  `json:"created"`
    Updated *time.Time `json:"updated,omitempty"`
}
```

### Custom User Models

For database-backed users, create an export model at `./.projectforge/export/models/user.json`:

```json
{
  "name": "user",
  "package": "user",
  "description": "A user of the system",
  "icon": "profile",
  "columns": [
    {
      "name": "id",
      "type": "uuid",
      "pk": true,
      "search": true
    },
    {
      "name": "name",
      "type": "string",
      "search": true,
      "tags": ["title"]
    },
    {
      "name": "email",
      "type": "string",
      "nullable": true,
      "search": true
    },
    {
      "name": "created",
      "type": "timestamp",
      "sqlDefault": "now()",
      "tags": ["created"]
    },
    {
      "name": "updated",
      "type": "timestamp",
      "nullable": true,
      "sqlDefault": "now()",
      "tags": ["updated"]
    }
  ]
}
```

**Required Fields**: `id` and `name` are required for all user models.

## Dependencies

### Required Modules
- **`core`** - Base infrastructure and utilities

### Optional Modules
- **`database`** - For database-backed user storage
- **`export`** - For custom user model generation
- **`oauth`** - For authentication integration
- **`audit`** - For user activity tracking

## Usage Examples

### Basic User Operations

```go
// Create a new user
user := &user.User{
    ID:      uuid.New(),
    Name:    "John Doe",
    Created: util.TimeCurrent(),
}

// Save user (implementation depends on storage backend)
err := userService.Save(ctx, user)
```

### Session Integration

```go
// Get current user from session
currentUser := session.GetUser(ps)
if currentUser == nil {
    // Handle anonymous user
}

// Update user in session
session.SetUser(ps, updatedUser)
```

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/user
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [OAuth Module](oauth.md) - Authentication integration
- [Database Module](database.md) - Database storage backend
- [Export Module](export.md) - Custom model generation
- [Audit Module](audit.md) - User activity tracking
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
