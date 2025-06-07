# Rich Editor

The **Rich Editor** module provides an advanced in-browser editing experience for structured data with progressive enhancement. It offers a sophisticated table-based editor that gracefully falls back to raw JSON editing when JavaScript is disabled.

## Overview

This module enables users to edit complex data structures through an intuitive table interface, with seamless toggling between visual editing and raw JSON modes. It's designed with **progressive enhancement** principles, ensuring full functionality even without JavaScript.

### Key Features

- **Visual Table Editor**: Edit structured data in an intuitive table format
- **Dual Mode Operation**: Toggle between visual editor and raw JSON
- **Progressive Enhancement**: Works without JavaScript, enhanced with it
- **Type-Aware Fields**: Intelligent input types based on data structure
- **Live Data Sync**: Real-time synchronization between editor and underlying JSON
- **Flexible Integration**: Multiple template variants for different UI contexts

## Architecture

### Client-Side (TypeScript)

The module is primarily implemented in TypeScript with five core files:

- **`editor.ts`** - Main editor initialization and state management
- **`editortypes.ts`** - Type definitions and utility functions
- **`editorfield.ts`** - Field-specific input handling (string, boolean, integer, etc.)
- **`editortable.ts`** - Dynamic table generation and row management
- **`editorobject.ts`** - Complex object and nested data handling

### Server-Side (Quicktemplate)

Server-side rendering through quicktemplate components:

- **`RichEdit.html`** - Main component templates with three variants

## Component Variants

### 1. Basic Rich Editor
```go
{%= RichEditor(key, id, title, columns, values, placeholder...) %}
```
Renders a standalone rich editor with toggle functionality.

### 2. Table Row Editor
```go
{%= RichEditorTable(key, id, title, columns, values, placeholder...) %}
```
Integrates the editor into an HTML table row format.

### 3. Card-Based Editor
```go
{%= RichEditorCard(key, id, title, ps, headerExtra, icon, columns, values, placeholder...) %}
```
Presents the editor within a card layout with header controls and icons.

## Data Structure

### Column Definition
```typescript
type Column = {
  key: string        // Field identifier
  title: string      // Display name
  description?: string // Optional help text
  type?: string      // Data type hint
}
```

### Editor State
```typescript
type Editor = {
  key: string                           // Unique editor identifier
  title: string                         // Editor display title
  columns: Column[]                     // Field definitions
  textarea: HTMLTextAreaElement         // Underlying JSON storage
  rows: { [key: string]: unknown; }[]   // Current data rows
  table?: Element                       // Generated table element
}
```

## Supported Field Types

The editor intelligently handles multiple data types:

- **string**: Text input fields with validation
- **bool**: Radio button pairs (True/False)
- **int**: Number inputs with parsing
- **float**: Decimal number inputs
- **enum**: Dropdown selection from predefined values
- **list**: Array handling with nested type support
- **object**: Complex nested data structures

## Usage Examples

### Basic Implementation
```html
<div class="rich-editor" 
     data-key="users" 
     data-title="User Management" 
     data-columns='[{"key":"name","title":"Name","type":"string"},{"key":"active","title":"Active","type":"bool"}]'>
  <textarea name="users">[{"name":"John","active":true}]</textarea>
</div>
```

### With Go Template Data
```go
{% code 
  columns := []*util.FieldDesc{
    {Key: "name", Title: "Full Name", Type: "string"},
    {Key: "email", Title: "Email Address", Type: "string"},
    {Key: "admin", Title: "Administrator", Type: "bool"},
  }
  users := []any{
    map[string]any{"name": "John Doe", "email": "john@example.com", "admin": false},
  }
%}

{%= RichEditorCard("users", "user-editor", "User Management", ps, "", "user", columns, users, "Manage application users") %}
```

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

## Integration

### Dependencies
- **TypeScript**: Client-side functionality
- **DOM utilities**: Provided by core module
- **Quicktemplate**: Server-side rendering

### Initialization
The editor automatically initializes on page load:

```typescript
// Automatic initialization
editorInit(); // Finds all .rich-editor elements

// Manual initialization
createEditor(element); // Initialize specific element
```

### CSS Classes
- `.rich-editor` - Main container class
- `.toggle-editor-{key}` - Toggle button selector
- `.rich-editor table` - Generated table styling

## Advanced Features

### Custom Type Handling
The module supports complex type definitions:

```typescript
// Simple type
type: "string"

// Complex type with metadata
type: {
  k: "enum",
  t: { ref: "user-roles" }
}

// List type
type: {
  k: "list", 
  t: { v: "string" }
}
```

### Dynamic Row Management
- Add new rows with default values
- Remove existing rows with confirmation
- Reorder rows (where supported)
- Bulk operations on selected rows

## Configuration

The module requires no additional configuration beyond the column definitions passed to the template functions. All behavior is determined by:

1. **Column Metadata**: Type information and field properties
2. **Initial Data**: Starting values for the editor
3. **Container Attributes**: Key, title, and column definitions

## Browser Compatibility

- **Modern Browsers**: Full TypeScript functionality
- **Legacy Browsers**: Graceful degradation to textarea
- **No JavaScript**: Complete fallback functionality
- **Screen Readers**: Accessible form controls

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/richedit
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete framework documentation
- [Quicktemplate Guide](https://github.com/valyala/quicktemplate) - Template syntax reference
- [TypeScript Documentation](https://www.typescriptlang.org/) - Language reference
