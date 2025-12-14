# JSON Schema

The **`jsonschema`** module provides comprehensive [JSON Schema](https://json-schema.org/) support for your application.
It offers a complete Go representation of JSON Schema Draft 2020-12 with validation utilities and collection management.

## Overview

This module provides:

- **Schema Representation**: Complete JSON Schema Draft 2020-12 specification compliance
- **Validation**: Schema structure validation and required field checking
- **Collection Management**: Organize and manage multiple schemas with ID-based lookup
- **JSON Utilities**: Parsing and serialization support for schema files

## Key Features

### JSON Schema Compliance
- Full JSON Schema Draft 2020-12 specification support
- Core vocabulary implementation (type, const, enum, etc.)
- Validation keywords (minimum, maximum, pattern, format, etc.)
- Boolean logic operations (allOf, anyOf, oneOf, not)
- Object and array validation (properties, items, additionalProperties, etc.)

### Schema Management
- Schema collection organization
- ID-based schema lookup and retrieval
- Validation for schema structure integrity
- Support for schema references and definitions

### Developer Experience
- Type-safe Go structs for all schema components
- JSON marshaling/unmarshaling support
- Clear validation error reporting
- Easy integration with existing Go applications

## Package Structure

### Core Components

- **`schema.go`** - Main Schema struct with complete JSON Schema specification
  - Schema metadata ($id, $schema, title, description)
  - Type definitions and constraints
  - Validation keywords and boolean logic
  - Object and array schema definitions

- **`collection.go`** - Schema collection management
  - Multiple schema organization
  - ID-based schema lookup
  - Collection validation utilities

- **`validate.go`** - Schema validation logic
  - Schema structure validation
  - Required field checking
  - Constraint validation

- **`json.go`** - JSON parsing and serialization
  - Schema file loading
  - JSON marshaling/unmarshaling
  - Error handling for malformed schemas

## Usage Examples

### Basic Schema Creation

```go
schema := &jsonschema.Schema{
    Type:        "object",
    Title:       "User Schema",
    Description: "Schema for user objects",
    Properties: map[string]*jsonschema.Schema{
        "name": {
            Type: "string",
            MinLength: util.IntPtr(1),
        },
        "age": {
            Type:    "integer",
            Minimum: util.IntPtr(0),
        },
    },
    Required: []string{"name"},
}
```

### Schema Collection Management

```go
collection := jsonschema.NewCollection()
collection.Add("user", userSchema)
collection.Add("product", productSchema)

userSchema := collection.Get("user")
```

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/jsonschema
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [JSON Schema Specification](https://json-schema.org/) - Official JSON Schema documentation
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
