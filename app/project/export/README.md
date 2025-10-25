# Export

This package provides the functionality to export Golang and TypeScript code based on the files in the `.projectforge/export` directory accessible from the project root.

### Key Terms

`Model`: A data structure that represents a full object and its fields. Models can generate a full database-backed service, a rich UI, and more. Defined in `app/lib/metamodel/model/model.go`

`Enum`: A data structure that represents a set of named constants. Enums can be used to define the possible values for a field in a model, and are backed by rich Go/TS objects. Defined in `app/lib/metamodel/enum/enum.go`

`Event`: A data structure an event and its fields. Basically a lightweight model, these aren't persisted (automatically) in a database. Defined in `app/lib/metamodel/model/event.go`

`Group`: A data structure that represents a group that models can belong to. Groups can be used to organize models into a hierarchy. Defined in `app/lib/metamodel/model/group.go`

### Export Directory Structure

The `export` directory can contain the following:

- `enum`: Directory of exported enums, each a JSON file that serializes to an `Enum`
- `model`: Directory of exported models, each a JSON file that serializes to a `Model`
- `event`: Directory of exported events, each a JSON file that serializes to an `Event`
- `groups.json`: A JSON file that contains an array of groups that models can belong to. They conform to `Group`.
- `types.json`: A JSON file that contains an array of extra types that can be used in the export. They conform to `Model`, but are not actually models used in the app

## Exported Files

### Enumerations

There are two types of enumerations. Simple enums are represented as a `string` alias, while complex enums are represented as a struct.

Both serialize to a string for JSON and friends, and provide utility functions.

#### "Simple" Enums

For the following configuration file, saved as `.projectforge/export/enums/widget.json`:

```json
{
  "name": "widget",
  "package": "inventory",
  "description": "Widget",
  "icon": "star",
  "values": ["a", "b", "c", "d"]
}
```

The following code will be generated as `app/inventory/widget.go`:

```go
package inventory

var AllWidgets = Widgets{WidgetA, WidgetB, WidgetC, WidgetD}

type (
	Widget  string
	Widgets []Widget
)

const (
	WidgetA Widget = "a"
	WidgetB Widget = "b"
	WidgetC Widget = "c"
	WidgetD Widget = "d"
)
```

#### "Complex" Enums

As soon as at least one enum value carries metadata (name, description, icon, default flag, or extras), an enum is considered "complex".

For the following configuration file, saved as `.projectforge/export/enums/box.json`:

```json
{
  "name": "box",
  "package": "inventory",
  "description": "A complex enum",
  "icon": "star",
  "values": [
    { "key": "small", "name": "Small Box", "extra": { "flaps": 1 } },
    { "key": "large", "name": "Large Box", "extra": { "flaps": 4 } }
  ]
}
```

...the following files will be emitted to `app/inventory/box*.go`:

- Struct definition `Box` with the fixed fields `Key`, `Name`, `Description`, `Icon`, plus one field per field defined in `extra`. Extra field types are inferred from sample data, can be overridden with `config.type:<field>`.
- For each value, a package-level variable like `BoxSmall` (a `*Box`) is declared, collecting all configured fields.
- Type alias `Boxes`, defined as `[]*Box`, and variable `AllBoxes`, aggregating the variables so the enum can be iterated.
- Collection alias `type Boxes []Box` with helpers: `Keys()`, `Strings()`, `Help() string` (comma-separated options), `Random()`.
- Lookup functions on `Boxes`:
  - `Get(key string, logger util.Logger)` compares case-insensitively against the `Key` field.
  - `GetByName(name string, logger util.Logger)` mirrors the same logic for the `Name` field.
  - For each extra field, a helper is emitted, such as `GetByFlaps(value int, logger util.Logger) *Box` if values are unique; otherwise `GetByFlaps(value) Boxes`.
- `func BoxParse(logger util.Logger, keys ...string) *Box` maps raw strings to the appropriate enum value by `Key`. Empty input returns `nil`.

Every lookup logs a warning through the supplied logger when no match is found. If an enum value is marked as the default, that value is returned on missing or blank input; otherwise a sentinel `{Key: "_error", Name: "error: ..."} ` is produced so callers can detect configuration issues.

### Events

Events are used to represent ephemeral structs that are (usually) not stored in the database. They're useful as WebSocket messages, event-driven pipelines, and really any time you need a rich struct in Go or TypeScript.

For the following configuration file, saved as `.projectforge/export/events/config.json`:

```json
{
  "name": "config",
  "package": "state",
  "group": ["game"],
  "description": "The configuration for this particular game state",
  "icon": "star",
  "columns": [
    { "name": "written_language", "type": "string" },
    { "name": "verbose_logging", "type": "bool" },
    { "name": "extra", "type": "map" }
  ]
}
```

...the following will be emitted to `app/game/state/config.go`:

- The `Config` struct, with correctly-typed fields, in this case `string`, `bool`, `util.ValueMap`.
- Two constructors, `NewConfig` and `ConfigFromMap`, and `RandomConfig` which is useful for testing.
- Transformation methods `ToMap`, `ToOrderedMap`, `Strings`, `ToCSV`, `ToData`.
- A `Clone` method to create a deep copy of the struct.
- A `Diff` method to compare two instances and return a list of differences.
- Type Alias `Configs []*Config`, with methods `Clone`, `ToMaps`, `ToOrderedMaps`, and accessors by every unique or indexed field.

### Models

Models are used to represent persistent structs that have a primary key (composite keys supported).
They are often stored in the database, and UI and managed service are required. They're useful as the foundation for your application's data model.

For the following configuration file, saved as `.projectforge/export/models/config.json`:

```json
{
  "name": "config",
  "package": "state",
  "description": "The configuration for this particular game state",
  "icon": "star",
  "columns": [
    { "name": "id", "type": "uuid", "pk": true },
    { "name": "written_language", "type": "string" },
    { "name": "verbose_logging", "type": "bool" },
    { "name": "extra", "type": "map" }
  ]
}
```

...several file will be emitted to `app/state`:

- `config.go`, containing:
  - The `Config` struct, with correctly-typed fields, in this case `uuid.UUID`, `string`, `bool`, `util.ValueMap`.
  - Two constructors, `NewConfig` and `ConfigFromMap`, and `RandomConfig` which is useful for testing.
  - Transformation methods `String`, `Strings`, `ToCSV`, `ToData`.
  - A `Clone` method to create a deep copy of the struct.

- `configdiff.go`
  - A `Diff` method to compare two instances and return a list of differences.

- `configmap.go`
  - Transformation methods `ToMap`, `FromMap`, `ToOrderedMap`.

- `configs.go`
  - Type Alias `Configs []*Config`, with methods `Clone`, `ToMaps`, `ToOrderedMaps`, and accessors by every unique or indexed field.
