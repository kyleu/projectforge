# Expression

The **`expression`** module provides Common Expression Language (CEL) integration for [Project Forge](https://projectforge.dev) applications. It enables runtime evaluation of arbitrary expressions with high performance and caching capabilities.

## Overview

This module provides:

- **CEL Engine**: Google's Common Expression Language runtime for safe expression evaluation
- **Expression Caching**: Performance-optimized caching with compile-once-run-many pattern
- **Type-Safe Evaluation**: Compile-time checking and runtime type validation
- **Boolean Result Checking**: Utilities for evaluating conditional expressions

## Key Features

### Performance
- Compile-once-run-many pattern with automatic caching
- Thread-safe expression compilation and evaluation
- Sub-millisecond evaluation times for cached expressions
- Zero-allocation result checking where possible

### Safety
- Sandboxed expression execution via CEL runtime
- Compile-time expression validation
- Type-safe parameter passing
- Error handling with detailed context

### Flexibility
- Support for complex expressions with variables
- Custom CEL environment options
- Extensible expression parameter types
- Boolean and arbitrary result evaluation

## Package Structure

### Core Components

- **`lib/expression/model.go`** - Expression struct and compilation/execution logic
  - Expression definition with key, description, and pattern
  - Compile method for CEL AST generation and program creation
  - Run method for parameterized expression evaluation
  - Result type checking utilities

- **`lib/expression/check.go`** - Engine with caching and evaluation services
  - Thread-safe expression cache management
  - CEL environment configuration and initialization
  - High-level Check method for boolean expression evaluation
  - Compile method for pre-compilation and caching

## Usage Examples

### Basic Expression Evaluation

```go
// Create CEL engine
engine, err := expression.NewEngine()
if err != nil {
    return err
}

// Evaluate boolean expression
params := map[string]any{
    "age": 25,
    "name": "John",
}
result, err := engine.Check("age >= 18 && size(name) > 0", params, logger)
// result: true
```

### Pre-compiled Expressions

```go
// Compile expression for reuse
expr, err := engine.Compile("price * quantity > 100", logger)
if err != nil {
    return err
}

// Run with different parameters
params1 := map[string]any{"price": 10.50, "quantity": 12}
result1, duration, err := expr.Run(params1)
// result1: true, runs in microseconds due to pre-compilation
```

### Custom CEL Environment

```go
// Create engine with custom CEL options
engine, err := expression.NewEngine(
    cel.Declarations(
        decls.NewVar("user", decls.NewObjectType("User")),
    ),
)
```

## CEL Expression Language

This module uses Google's [Common Expression Language (CEL)](https://github.com/google/cel-spec), which provides:

### Supported Operators
- **Logical**: `&&`, `||`, `!`
- **Comparison**: `==`, `!=`, `<`, `<=`, `>`, `>=`
- **Arithmetic**: `+`, `-`, `*`, `/`, `%`
- **String**: `+` (concatenation), `matches()`, `contains()`
- **Collection**: `in`, `size()`, indexing with `[]`

### Built-in Functions
- **String**: `size()`, `startsWith()`, `endsWith()`, `matches()`
- **Collection**: `size()`, `in` operator
- **Type**: `type()`, conversion functions
- **Math**: Standard arithmetic operations

### Data Types
- **Primitives**: `bool`, `int`, `uint`, `double`, `string`, `bytes`
- **Collections**: `list`, `map`
- **Well-known types**: `timestamp`, `duration`

## Configuration

The module supports CEL environment customization via engine creation options:

```go
engine, err := expression.NewEngine(
    cel.Types(&CustomType{}),                    // Custom type definitions
    cel.Declarations(decls.NewVar(...)),         // Variable declarations
    cel.Functions(customFunctions...),           // Custom function definitions
)
```

## Performance Characteristics

- **Compilation**: One-time cost per unique expression pattern
- **Evaluation**: Sub-millisecond for most expressions
- **Memory**: Expressions cached indefinitely (consider memory bounds for long-running applications)
- **Concurrency**: Thread-safe with minimal contention via cache-first pattern

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/expression
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [CEL Specification](https://github.com/google/cel-spec) - Complete CEL language reference
- [CEL-Go Documentation](https://github.com/google/cel-go) - Go implementation details
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
