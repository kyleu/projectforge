# Settings

The **`settings`** module provides a file-backed configuration store for app-wide settings.
It persists a typed `Settings` struct to `settings.json` in the application's configuration directory and exposes a small service for reading and writing those values.

## Overview

This module gives you a simple, reliable way to store application settings that are:

- **File-backed**: Stored as JSON in the config directory for portability and version control.
- **Type-safe**: Accessed via a generated `Settings` struct rather than loose maps.
- **Admin-editable**: Exposed in the admin UI with a built-in editor.
- **Cache-friendly**: Read once, cached in memory, and cloned on access to avoid mutation bugs.

## Key Features

### Persistence and Safety
- JSON storage at `settings.json` with pretty formatting.
- Unknown keys are rejected when saving from a map, preventing silent misconfiguration.
- Simple clone/diff helpers for safe comparisons and debugging.

### Admin UI Integration
- `/admin` shows the current settings in a read-only table.
- `/admin/settings` offers a form-backed editor with field descriptions.

### Developer Convenience
- `ToMap`, `ToOrderedMap`, and `ToCSV` utilities for export or diagnostics.
- `SettingsFieldDescs` provides labels/types for UI rendering.

## Package Structure

### Core Logic

- **`app/lib/settings/settings.go`** - Settings struct and helpers
  - Field descriptors for admin UI
  - Map/CSV conversions and diffing
  - Generated once for customization

- **`app/lib/settings/service.go`** - Settings service
  - Cached `Get` and `Sync` operations
  - `Set` and `SetMap` with JSON persistence
  - Backed by the filesystem module

### Admin Views

- **`views/vadmin/SettingsEdit.html`** - Editable settings form
  - Uses `SettingsFieldDescs` to render inputs

## Usage Examples

### Read and Update Settings

```go
st := as.Services.Settings.Get()
st.ExampleBool = true
if err := as.Services.Settings.Set(st); err != nil {
    return err
}
```

### Update from a Map

```go
err := as.Services.Settings.SetMap(util.ValueMap{
    "exampleBool":   true,
    "exampleString": "updated",
})
if err != nil {
    return err
}
```

### Extend the Settings Schema

Update the generated settings file and field descriptors:

```go
// app/lib/settings/settings.go
type Settings struct {
    ExampleBool   bool   `json:"exampleBool,omitzero"`
    ExampleString string `json:"exampleString,omitzero"`
    NewFlag       int    `json:"newFlag,omitzero"`
}
```

Add a corresponding entry to `SettingsFieldDescs` so the admin UI can render it.

## Configuration

- **Location**: `settings.json` is stored in the app config directory, the same root used by the filesystem module.
- **Config directory**: Typically set via the `--config_dir` flag or the OS-specific default chosen at startup.

## Dependencies

### Required Modules
- **`filesystem`** - Provides the storage backend for `settings.json`

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/settings
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Filesystem Module](filesystem.md) - Core filesystem connectivity
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
