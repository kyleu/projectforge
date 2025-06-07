# Types

The **`types`** module provides a comprehensive type system for [Project Forge](https://projectforge.dev) applications. It offers a collection of structured data types with built-in validation, serialization, and type safety features.

## Overview

This module provides a unified type system for representing and working with common data types in Go applications. It focuses on:

- **Type Safety**: Strongly typed representations with compile-time guarantees
- **JSON Serialization**: Built-in support for JSON marshaling and unmarshaling
- **Validation**: Configurable constraints and validation rules
- **Extensibility**: Composable type system with interfaces for custom types

## Key Features

### Comprehensive Type Coverage
- **Primitive Types**: String, Int, Float, Bool, Byte, Char
- **Complex Types**: List, Map, Set, OrderedMap, ValueMap
- **Temporal Types**: Date, Time, Timestamp, TimestampZoned
- **Special Types**: UUID, JSON, XML, Enum, Reference, Option
- **Utility Types**: Nil, Any, Unknown, Wrapped, Error

### Advanced Features
- **Type Conversion**: Safe conversion between compatible types
- **Default Values**: Configurable default value generation
- **Sorting Support**: Built-in sortability indicators
- **Scalar Detection**: Distinction between scalar and composite types
- **Range Constraints**: Min/max validation for numeric and string types

### JSON Integration
- **Custom Marshaling**: Optimized JSON serialization for all types
- **Flexible Parsing**: Intelligent type detection from JSON values
- **Wrapped Types**: Support for complex nested type structures

## Package Structure

### Core Types

- **`type.go`** - Core type interface and utilities
  - `Type` interface definition
  - Type casting and conversion utilities
  - Common type checking functions

- **`string.go`** - String type with length and pattern validation
- **`int.go`** - Integer type with range constraints
- **`float.go`** - Floating-point type with precision controls
- **`bool.go`** - Boolean type representation

### Collection Types

- **`list.go`** - Ordered collections with type constraints
- **`map.go`** - Key-value mappings with typed keys and values
- **`set.go`** - Unique value collections
- **`orderedmap.go`** - Ordered key-value mappings
- **`valuemap.go`** - Specialized value mappings

### Temporal Types

- **`date.go`** - Date-only representations
- **`time.go`** - Time-only representations
- **`timestamp.go`** - Combined date and time
- **`timestampzoned.go`** - Timezone-aware timestamps

### Specialized Types

- **`uuid.go`** - UUID type with validation
- **`json.go`** - Raw JSON value handling
- **`xml.go`** - XML document representation
- **`enum.go`** - Enumeration type definitions
- **`option.go`** - Optional value handling

### Utility Types

- **`any.go`** - Dynamic type representation
- **`nil.go`** - Null value handling
- **`unknown.go`** - Unknown type representation
- **`wrapped.go`** - Type composition and nesting
- **`reference.go`** - Reference type for complex relationships

## Usage Examples

### Basic Type Usage

```go
// String type with constraints
stringType := &types.String{
    MinLength: 1,
    MaxLength: 100,
    Pattern:   "^[a-zA-Z0-9]+$",
}

// Integer type with range
intType := &types.Int{
    Min: 0,
    Max: 1000,
}

// Convert and validate values
value := stringType.From("hello")
if stringType.Sortable() {
    // Handle sortable type
}
```

### Collection Types

```go
// List of strings
listType := &types.List{
    Val: &types.String{MaxLength: 50},
}

// Map with string keys and integer values
mapType := &types.Map{
    K: &types.String{},
    V: &types.Int{Min: 0},
}
```

### Type Checking

```go
if types.IsString(someType) {
    // Handle string type
}

if types.IsList(someType) {
    listType := types.TypeAs[*types.List](someType)
    // Work with list type
}
```

## Configuration

The types module supports various configuration options through individual type properties:

- **String constraints**: MinLength, MaxLength, Pattern
- **Numeric ranges**: Min, Max values for Int and Float
- **Collection settings**: Value types for List, key/value types for Map
- **Temporal formats**: Custom formatting for date/time types

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/types
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
