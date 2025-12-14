# Edit Component Helpers

This directory contains quicktemplate functions for various form input components used throughout {{{ .Name }}}.
All functions are defined as quicktemplate components in `.html` files.

## AnyByType.html
{{{ if .HasModule "types" }}}
### Type Dispatch Functions
- `AnyByType(key string, id string, x any, t *types.Wrapped)` - Routes to appropriate edit component based on type
- `Default(key string, id string, x any, t types.Type)` - Fallback handler for unhandled types
{{{ end }}}
## Array.html (Selection Components)

### Select Dropdowns
- `Select(key string, id string, value string, opts []string, titles []string, indent int)` - Basic select dropdown
- `SelectVertical(key string, id string, title string, value string, opts []string, titles []string, indent int, help ...string)` - Vertical layout select
- `SelectTable(key string, id string, title string, value string, opts []string, titles []string, indent int, help ...string)` - Table row select

### Datalist Components
- `Datalist(key string, id string, value string, opts []string, titles []string, indent int, placeholder ...string)` - Input with datalist suggestions
- `DatalistVertical(key string, id string, title string, value string, opts []string, titles []string, indent int, help ...string)` - Vertical layout datalist
- `DatalistTable(key string, id string, title string, value string, opts []string, titles []string, indent int, help ...string)` - Table row datalist

### Radio Buttons
- `Radio(key string, value string, opts []string, titles []string, indent int)` - Basic radio button group
- `RadioVertical(key string, title string, value string, opts []string, titles []string, indent int, help ...string)` - Vertical layout radio group
- `RadioTable(key string, title string, value string, opts []string, titles []string, indent int, help ...string)` - Table row radio group

### Checkboxes
- `Checkbox(key string, values []string, opts []string, titles []string, indent int)` - Basic checkbox group
- `CheckboxVertical(key string, title string, values []string, opts []string, titles []string, indent int, help ...string)` - Vertical layout checkbox group
- `CheckboxTable(key string, title string, values []string, opts []string, titles []string, indent int, help ...string)` - Table row checkbox group

## Bool.html

### Boolean Input Components
- `Bool(key string, id string, x any, nullable bool)` - Radio buttons for true/false/null
- `BoolVertical(key string, title string, value bool, indent int, help ...string)` - Vertical layout boolean input
- `BoolTable(key string, title string, value bool, indent int, help ...string)` - Table row boolean input

## Color.html

### Color Input Components
- `Color(key string, id string, color string, placeholder ...string)` - HTML5 color picker input
- `ColorTable(key string, id string, title string, color string, indent int, help ...string)` - Table row color picker

## File.html

### File Upload Components
- `File(key string, id string, label string, value string, placeholder ...string)` - Single file upload
- `FileTable(key string, id string, title string, label string, value string)` - Table row single file upload
- `FileMultiple(key string, id string, label string, value string, placeholder ...string)` - Multiple file upload
- `FileMultipleTable(key string, id string, title string, label string, value string)` - Table row multiple file upload

## Float.html

### Float Number Input Components
- `Float(key string, id string, value any, placeholder ...string)` - Decimal number input
- `FloatVertical(key string, id string, title string, value float64, indent int, help ...string)` - Vertical layout float input
- `FloatTable(key string, id string, title string, value float64, indent int, help ...string)` - Table row float input

## Icon.html

### Icon Selection Components
- `IconPicker(key string, selected string, ps *cutil.PageState, indent int)` - Visual icon picker with radio buttons
- `IconPickerVertical(key string, title string, value string, ps *cutil.PageState, indent int)` - Vertical layout icon picker
- `Icons(key string, title string, value string, ps *cutil.PageState, indent int)` - Icon selection component
- `IconsTable(key string, title string, value string, ps *cutil.PageState, indent int, help ...string)` - Table row icon picker

## Int.html

### Integer Input Components
- `Int(key string, id string, value any, placeholder ...string)` - Integer number input
- `IntVertical(key string, id string, title string, value int, indent int, help ...string)` - Vertical layout integer input
- `IntTable(key string, id string, title string, value int, indent int, help ...string)` - Table row integer input

## Option.html

### Optional Type Components
- `Option(key string, id string, x any, t *types.Option)` - Handles optional/nullable types with null button
{{{ if .HasModule "richedit" }}}
## RichEdit.html

### Rich Editor Components
- `RichEditor(key string, id string, title string, columns util.FieldDescs, values []any, placeholder ...string)` - JSON-based rich text editor
- `RichEditorTable(key string, id string, title string, columns util.FieldDescs, values []any, placeholder ...string)` - Table row rich editor
- `RichEditorCard(key string, id string, title string, ps *cutil.PageState, headerExtra string, icon string, columns util.FieldDescs, values []any, placeholder ...string)` - Card-wrapped rich editor with toggle
{{{ end }}}
## Search.html

### Search Form Components
- `SearchForm(action string, fieldKey string, placeholder string, value string, ps *cutil.PageState)` - Search form with preserved query parameters

## String.html

### Text Input Components
- `String(key string, id string, value string, placeholder ...string)` - Basic text input
- `StringVertical(key string, id string, title string, value string, indent int, help ...string)` - Vertical layout text input
- `StringTable(key string, id string, title string, value string, indent int, help ...string)` - Table row text input

### Password Input Components
- `Password(key string, id string, value string, placeholder ...string)` - Password input
- `PasswordVertical(key string, id string, title string, value string, indent int, help ...string)` - Vertical layout password input
- `PasswordTable(key string, id string, title string, value string, indent int, help ...string)` - Table row password input

### Textarea Components
- `Textarea(key string, id string, rows int, value string, placeholder ...string)` - Multi-line text input
- `TextareaVertical(key string, id string, title string, rows int, value string, indent int, help ...string)` - Vertical layout textarea
- `TextareaTable(key string, id string, title string, rows int, value string, indent int, help ...string)` - Table row textarea

## TableEdit.html

### Table Editor Components
- `TableEditor(key string, columns util.FieldDescs, values util.ValueMap, action string, method string, title string)` - Complete table editor with form wrapper
- `TableEditorNoForm(key string, columns util.FieldDescs, values util.ValueMap, name string, value string, title string)` - Table editor without form wrapper
- `TableEditorNoTable(key string, columns util.FieldDescs, values util.ValueMap)` - Just the table rows without table wrapper

## Tags.html

### Tag Input Components
- `Tags(key string, id string, values []string, ps *cutil.PageState, placeholder ...string)` - Interactive tag editor with add/remove functionality
- `TagsVertical(key string, id string, title string, values []string, ps *cutil.PageState, indent int, help ...string)` - Vertical layout tag editor
- `TagsTable(key string, id string, title string, values []string, ps *cutil.PageState, indent int, help ...string)` - Table row tag editor

## Timestamp.html

### Timestamp Input Components
- `Timestamp(key string, id string, value *time.Time, placeholder ...string)` - Date/time picker input
- `TimestampVertical(key string, id string, title string, value *time.Time, indent int, help ...string)` - Vertical layout timestamp input
- `TimestampTable(key string, id string, title string, value *time.Time, indent int, help ...string)` - Table row timestamp input

### Date-only Input Components
- `TimestampDay(key string, id string, value *time.Time, placeholder ...string)` - Date-only picker
- `TimestampDayVertical(key string, id string, title string, value *time.Time, indent int, help ...string)` - Vertical layout date picker
- `TimestampDayTable(key string, id string, title string, value *time.Time, indent int, help ...string)` - Table row date picker

## UUID.html

### UUID Input Components
- `UUID(key string, id string, value *uuid.UUID, placeholder ...string)` - UUID text input (delegates to String)
- `UUIDVertical(key string, id string, title string, value *uuid.UUID, indent int, help ...string)` - Vertical layout UUID input
- `UUIDTable(key string, id string, title string, value *uuid.UUID, indent int, help ...string)` - Table row UUID input

## Usage Patterns

### Layout Variants
Most input components come in three layout variants:
- **Basic**: Just the input element (e.g., `String`, `Int`, `Float`)
- **Vertical**: Input with label above in a vertical layout (e.g., `StringVertical`, `IntVertical`)
- **Table**: Input formatted as a table row with label in left column (e.g., `StringTable`, `IntTable`)

### Common Parameters
- `key`: The form field name attribute
- `id`: The HTML id attribute (optional for basic variants)
- `title`: The label text for vertical and table variants
- `indent`: Indentation level for proper HTML formatting
- `help`: Optional help text that becomes a tooltip or description
- `placeholder`: Optional placeholder text for inputs

### Integration
These components are designed to work with the MVC architecture and integrate with:
- `cutil.PageState` for managing page state and assets
- `util.FieldDescs` for dynamic form field descriptions
- The application's CSS framework for consistent styling
