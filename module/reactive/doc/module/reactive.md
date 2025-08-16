# Reactive Values

The **`reactive`** module provides thread-safe reactive values with observer pattern support.

## Overview

This module is designed for applications that need **reactive values** and provides:

- Thread-safe reactive values with generic type support
- Observer pattern implementation for value change notifications
- Computed values that automatically update based on functions
- Simple API for subscribing to value changes

## Key Features

### Reactive Values
- Generic `Value[T]` type that holds any type of data
- Thread-safe Get/Set operations using sync.RWMutex
- Observer pattern with Subscribe method for value change notifications
- Automatic notification of all observers when value changes

### Computed Values
- `Computed[T]` type that extends reactive values
- Values computed from functions that can be recomputed on demand
- Built on top of reactive values, inheriting all observer functionality
- Manual recomputation via Recompute() method

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/reactive
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
