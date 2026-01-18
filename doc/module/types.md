# Types

The **`types`** module provides a shared type catalog and helpers used across generated apps. It defines type metadata, JSON-friendly type definitions, and view helpers for rendering or editing values by type.

## Overview

- `types.Type` describes a value's key, sortability, scalar-ness, conversion (`From`), and default value (`Default`).
- `types.Wrapped` wraps concrete types for JSON round-tripping and convenience constructors (`NewString`, `NewList`, and friends).
- Helper utilities: `TypeAs`, `IsString`/`IsBool`/`IsInt`/`IsList`, `Bits`, and `FromReflect`.
- View components: `views/components/view/AnyByType` and `views/components/edit/AnyByType` render values and inputs based on a `types.Wrapped`.

## Type Catalog

- **Primitive**: `any`, `bit`, `bool`, `byte`, `char`, `int`, `float`, `numeric`, `string`
- **Temporal**: `date`, `time`, `timestamp`, `timestampZoned`
- **Identifiers / references**: `uuid`, `enum`, `enumValue`, `reference`
- **Structured**: `list`, `map`, `orderedMap`, `set`, `valueMap`, `numericMap`, `option`, `range`, `union`, `method`
- **Special**: `json`, `xml`, `nil`, `unknown`, `error`

Notes:
- `valueMap` is a string-keyed map of arbitrary values.
- `numeric` and `numericMap` are intended to pair with the `numeric` module's implementations.

## Usage

### Define types

```go
str := types.NewString()
name := types.NewStringArgs(1, 80, "^[a-zA-Z0-9_-]+$")
scores := types.NewList(types.NewInt(64))
valueMap := types.NewValueMap()
```

### JSON definitions

```go
_ = util.ToJSONCompact(types.NewString()) // "string"
_ = util.ToJSONCompact(types.NewList(types.NewString())) // {"k":"list","t":{"v":"string"}}

_, _ = util.FromJSONObj[*types.Wrapped]([]byte(`"int"`))
```

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/types
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Numeric Module](numeric.md) - Numeric runtime types
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
