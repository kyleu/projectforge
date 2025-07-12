# Metamodel

The **`metamodel`** module provides a comprehensive Go representation of database models, enumerations, and type definitions for [Project Forge](https://projectforge.dev) applications. It serves as the foundation for code generation, database schema management, and type-safe data modeling.

## Overview

This module enables applications to define and manipulate data structures as first-class objects, providing:

- **Model Definitions**: Complete database table representations with columns, relations, and constraints
- **Enumeration Types**: Type-safe enum definitions with values and configuration
- **Schema Management**: Database schema representation and validation
- **Code Generation Support**: Template-ready structures for generating application code

## Key Features

### Model Management
- Complete table/model definitions with metadata
- Column specifications with types, constraints, and validation
- Relationship definitions (foreign keys, associations)
- Index management for performance optimization
- Seed data specification for testing and development

### Type Safety
- Strongly-typed column definitions
- Validation rules and constraints
- Import dependency tracking
- Configuration value mapping

### Code Generation
- Template-ready data structures
- Package and naming convention support
- Extensible configuration system
- Cross-reference relationship handling

## Package Structure

### Core Components

- **`metamodel/model/`** - Database model representations
  - [`Model`](file:///Users/kyle/kyleu/projectforge/module/metamodel/app/lib/metamodel/model/model.go#L14-L39) - Complete table definition with columns, relations, indexes
  - [`Column`](file:///Users/kyle/kyleu/projectforge/module/metamodel/app/lib/metamodel/model/column.go#L32-L54) - Column specification with type, constraints, validation
  - [`Relation`](file:///Users/kyle/kyleu/projectforge/module/metamodel/app/lib/metamodel/model/relation.go) - Foreign key and association definitions
  - [`Index`](file:///Users/kyle/kyleu/projectforge/module/metamodel/app/lib/metamodel/model/index.go) - Database index specifications

- **`metamodel/enum/`** - Enumeration type definitions
  - [`Enum`](file:///Users/kyle/kyleu/projectforge/module/metamodel/app/lib/metamodel/enum/enum.go#L15-L29) - Enumeration with values and configuration
  - [`Value`](file:///Users/kyle/kyleu/projectforge/module/metamodel/app/lib/metamodel/enum/value.go) - Individual enum value definitions

### Support Utilities

- **Model Validation** - Schema validation and constraint checking
- **Type Mapping** - Go type system integration
- **Import Management** - Dependency resolution and code generation
- **Example Generation** - Sample data and test fixture creation

## Usage Examples

### Defining a Model

```go
model := &model.Model{
    Name:        "User",
    Package:     "user",
    Description: "User account information",
    Columns: model.Columns{
        {Name: "id", Type: types.NewUUID(), PK: true},
        {Name: "email", Type: types.NewString(), Nullable: false},
        {Name: "created", Type: types.NewTimestamp(), Nullable: false},
    },
    Indexes: model.Indexes{
        {Name: "user_email_idx", Columns: []string{"email"}, Unique: true},
    },
}
```

### Defining an Enumeration

```go
enum := &enum.Enum{
    Name:        "Status",
    Description: "User account status",
    Values: enum.Values{
        {Key: "active", Name: "Active"},
        {Key: "inactive", Name: "Inactive"},
        {Key: "suspended", Name: "Suspended"},
    },
}
```

## Configuration

The metamodel module integrates with Project Forge's configuration system and supports:

- **Schema Validation**: Automatic validation of model definitions
- **Naming Conventions**: Configurable naming patterns for generated code
- **Type Mapping**: Custom type mappings for different databases
- **Code Generation**: Template customization and output control

## Dependencies

- **Required**: [`types`](types.md) - Basic type system support

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/metamodel
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Types Module](types.md) - Basic type system foundation
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
