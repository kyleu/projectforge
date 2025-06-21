# Rich Editor

The **Rich Editor** module provides an advanced in-browser editing experience for structured data with progressive enhancement. It offers a sophisticated table-based editor that gracefully falls back to raw JSON editing when JavaScript is disabled.

### Key Features

- **Visual Table Editor**: Edit structured data in an intuitive table format
- **Dual Mode Operation**: Toggle between visual editor and raw JSON
- **Progressive Enhancement**: Works without JavaScript, enhanced with it
- **Type-Aware Fields**: Intelligent input types based on data structure
- **Live Data Sync**: Real-time synchronization between editor and underlying JSON
- **Flexible Integration**: Multiple template variants for different UI contexts

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

{%= components.RichEditor("users", "user-editor", "User Management", ps, "", "user", columns, users, "Manage application users") %}
```

### 1. Basic Rich Editor
```go
{%= components.RichEditor(key, id, title, columns, values, placeholder...) %}
```
Renders a standalone rich editor with toggle functionality.

### 2. Table Row Editor
```go
{%= components.RichEditorTable(key, id, title, columns, values, placeholder...) %}
```
Integrates the editor into an HTML table row format.

### 3. Card-Based Editor
```go
{%= components.RichEditorCard(key, id, title, ps, headerExtra, icon, columns, values, placeholder...) %}
```
Presents the editor within a card layout with header controls and icons.


## Progressive Enhancement

### Without JavaScript
- Displays as a standard textarea with JSON content
- Users can edit raw JSON directly
- Form submission works normally
- No visual editor interface

### With JavaScript Enabled
- Automatically converts to table-based editor
- Toggle button appears for switching modes
- Real-time validation and type checking
- Enhanced user experience with visual editing
