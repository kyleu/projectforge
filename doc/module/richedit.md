# Rich Editor

The **`richedit`** module provides an in-browser editor for structured JSON data with progressive enhancement. It renders a table-based editor when JavaScript is available and falls back to raw JSON editing when it is not.

## Overview

This module adds:
- Quicktemplate components in `views/components/edit/RichEdit.html`
- Client-side table editor with modal row editing
- Styles for `.rich-editor` and toggle controls

## Key Features

- **Visual Table Editor**: Click rows to edit data in a modal form
- **Dual Mode Operation**: Toggle between table view and raw JSON
- **Progressive Enhancement**: Works with or without JavaScript
- **Type-Aware Fields**: Input widgets driven by `util.FieldDesc` metadata
- **Flexible Embedding**: Standalone, table-row, and card variants

## Template Usage

```go
{% code
  columns := util.FieldDescs{
    {Key: "name", Title: "Full Name", Type: "string"},
    {Key: "email", Title: "Email Address", Type: "string"},
    {Key: "admin", Title: "Administrator", Type: "bool"},
  }
  users := []any{
    map[string]any{"name": "John Doe", "email": "john@example.com", "admin": false},
  }
%}

{%= components.RichEditorCard("users", "user-editor", "User Management", ps, "", "user", columns, users, "Manage application users") %}
```

### Component Variants

1. **Basic Rich Editor**
```go
{%= components.RichEditor(key, id, title, columns, values, placeholder...) %}
```

2. **Table Row Editor**
```go
{%= components.RichEditorTable(key, id, title, columns, values, placeholder...) %}
```

3. **Card-Based Editor**
```go
{%= components.RichEditorCard(key, id, title, ps, headerExtra, icon, columns, values, placeholder...) %}
```

### Custom Toggle Placement

If you want the toggle button outside of `RichEditorCard`, add a button with a matching class name:

```html
<button type="button" class="toggle-editor toggle-editor-users">Editor</button>
```

The script will change the label to "Raw JSON" when the table view is active.

## Data Model

- `columns` is a `util.FieldDescs` slice; each entry describes a column.
- `values` must be JSON-serializable and should be an array of objects.
- The editor stores the canonical JSON in the hidden textarea named `key`.

### Column Types

Input widgets are chosen by `FieldDesc.Type`:

- `bool` -> radio buttons
- `int` -> number input
- `type` -> JSON textarea for structured values
- any other value -> text input

Table rendering also respects these types:

- `code` or `json` -> pretty-printed JSON in the table
- `type` -> human-readable type string
- default -> `unknownToString` rendering

## Progressive Enhancement

### Without JavaScript
- Displays as a standard textarea with JSON content
- Users can edit raw JSON directly
- Form submission works normally
- Toggle controls stay hidden

### With JavaScript Enabled
- Textarea is replaced with a table view
- Clicking rows opens a modal editor
- Toggle button switches between table and raw JSON
- Textarea stays in sync with the edited rows

## Configuration

This module has no environment variables. Configure behavior through template parameters:

- `columns` defines the field metadata and types
- `values` provides the initial row data
- `placeholder...` supplies help text (also used as card title tooltip)
- optional toggle button with `toggle-editor` + `toggle-editor-<key>` classes

## CLI / URLs

This module does not add CLI commands or server routes. It is a UI component with client-side behavior only.

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/richedit
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
