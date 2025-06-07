# Numeric

The **`numeric`** module for [Project Forge](https://projectforge.dev) provides high-precision arithmetic operations for handling arbitrarily large numbers that exceed the limits of standard floating-point types. It implements mantissa-exponent representation for both Go and TypeScript environments.

## Overview

This module provides:

- **Large Number Support**: Handle numbers beyond standard 64-bit floating-point precision
- **Cross-Platform**: Identical API and behavior in both Go and TypeScript
- **Performance Optimized**: Pool-based memory management and efficient algorithms
- **Comprehensive Operations**: Full arithmetic, comparison, and formatting capabilities
- **Multiple Formats**: Standard, scientific, engineering, and English notation

## Key Features

### Precision
- Mantissa-exponent representation for arbitrary precision
- Support for numbers beyond IEEE 754 limits
- Consistent behavior across platforms

### Operations
- Complete arithmetic operations (add, subtract, multiply, divide, power)
- Comparison operations (equals, greater than, less than)
- Mathematical functions (absolute value, sign, normalization)
- Format conversions (string, scientific notation, engineering notation)

### Performance
- Object pooling in TypeScript for memory efficiency
- Zero-allocation utilities where possible
- Optimized for common mathematical operations

## Package Structure

### Go Implementation

- **`app/util/numeric/`** - Core Go implementation
  - `Numeric` struct with mantissa/exponent fields
  - Arithmetic and comparison operations
  - String formatting and parsing
  - JSON marshaling/unmarshaling support

### TypeScript Implementation

- **`client/src/numeric/`** - Full-featured TypeScript class
  - Pool-based memory management
  - Complete arithmetic operations
  - Multiple formatting options
  - Utility functions and constants

## Usage Examples

### Go
```go
import "myapp/app/util/numeric"

// Create numbers
a := numeric.FromString("123456789012345678901234567890")
b := numeric.FromFloat(1.23e15)

// Arithmetic operations
sum := a.Add(b)
product := a.Multiply(b)

// Formatting
fmt.Println(sum.String())          // Standard notation
fmt.Println(sum.Scientific())      // Scientific notation
fmt.Println(sum.Engineering())     // Engineering notation
```

### TypeScript
```typescript
import { Numeric } from './numeric/numeric';

// Create numbers
const a = new Numeric("123456789012345678901234567890");
const b = Numeric.fromFloat(1.23e15);

// Arithmetic operations
const sum = a.add(b);
const product = a.multiply(b);

// Formatting
console.log(sum.toString());       // Standard notation
console.log(sum.toScientific());   // Scientific notation
console.log(sum.toEngineering());  // Engineering notation
```

## Core Types

### Go Types
- **`Numeric`** - Main numeric type with mantissa and exponent
- **`mantissa`** - float64 for the significant digits
- **`exponent`** - int64 for the power of 10

### TypeScript Types
- **`Numeric`** - Main class with pooling support
- **`MantissaExponent`** - Interface for mantissa/exponent pairs
- **`NumericSource`** - Union type for input sources (string, number, or MantissaExponent)

## Configuration

No additional configuration required. The module works out of the box with sensible defaults for precision and formatting.

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/numeric
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)  
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
