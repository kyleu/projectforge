# View Component Helpers

This directory contains quicktemplate component functions for rendering various data types and UI elements.

## Function Reference

### Any.html

- `{% func Any(x any, ps *cutil.PageState) %}` - Renders any value by type-switching to appropriate handlers

### AnyByType.html

- `{% func AnyByType(x any, t *types.Wrapped, ps *cutil.PageState) %}` - Renders any value based on explicit type information
- `{% func Default(x any, t string, ps *cutil.PageState) %}` - Default handler for unhandled types
- `{% func Type(v types.Type) %}` - Renders type information
- `{% func Option(x any, t *types.Option, ps *cutil.PageState) %}` - Renders optional values with null toggle

### Bool.html

- `{% func Bool(b bool) %}` - Renders boolean as "true"/"false" text
- `{% func BoolIcon(b bool, size int, cls string, ps *cutil.PageState, titles ...string) %}` - Renders boolean as check/times icon

### Color.html

- `{% func Color(clr string, cls string, ps *cutil.PageState) %}` - Renders color value with background preview

### Diff.html

- `{% func Diffs(value util.Diffs) %}` - Renders diff table showing path/old/new changes
- `{% func DiffsSet(key string, value util.DiffsSet, limit int, ps *cutil.PageState) %}` - Renders grouped diffs as accordion

### Float.html

- `{% func Float(f any) %}` - Renders float value
- `{% func FloatArray(value []any) %}` - Renders array of floats

### Int.html

- `{% func Int(i any) %}` - Renders integer value
- `{% func IntArray(value []any) %}` - Renders array of integers

### Map.html

- `{% func Map(preserveWhitespace bool, m util.ValueMap, ps *cutil.PageState) %}` - Renders key-value map as table
- `{% func MapKeys(m util.ValueMap) %}` - Renders map keys as tags
- `{% func MapArray(preserveWhitespace bool, ps *cutil.PageState, maps ...util.ValueMap) %}` - Renders array of maps as table
- `{% func OrderedMap(preserveWhitespace bool, m *util.OrderedMap[any], ps *cutil.PageState) %}` - Renders ordered map as table
- `{% func OrderedMapArray(preserveWhitespace bool, ps *cutil.PageState, maps ...*util.OrderedMap[any]) %}` - Renders array of ordered maps as table

### Package.html

- `{% func Package(v util.Pkg) %}` - Renders package path segments

### String.html

- `{% func String(value string, classes ...string) %}` - Renders string with optional CSS classes
- `{% func StringRich(value string, code bool, maxLength int, classes ...string) %}` - Renders string with code formatting and length limit
- `{% func StringArray(value []string) %}` - Renders string array with overflow handling
- `{% func FormatLang(v string, ext string) %}` - Renders formatted code with syntax highlighting

### Tags.html

- `{% func Tags(values []string, titles []string, url ...string) %}` - Renders string array as tag elements with optional links

### Timestamp.html

- `{% func Timestamp(value *time.Time) %}` - Renders timestamp with verbose tooltip
- `{% func TimestampMillis(value *time.Time) %}` - Renders timestamp with millisecond precision
- `{% func TimestampRelative(value *time.Time, static bool) %}` - Renders relative time ("2 hours ago")
- `{% func TimestampDay(value *time.Time) %}` - Renders date only (YYYY-MM-DD)
- `{% func DurationSeconds(seconds float64) %}` - Renders duration from seconds value

### URL.html

- `{% func URL(u any, content string, includeExternalIcon bool, ps *cutil.PageState) %}` - Renders clickable URL link
- `{% func CodeLink(path string, title string, ps *cutil.PageState) %}` - Renders link to source code

### UUID.html

- `{% func UUID(value *uuid.UUID) %}` - Renders UUID as string

## Usage Notes

- Most functions accept `ps *cutil.PageState` for context and theme information
- Functions starting with capital letters are exported quicktemplate functions
- `preserveWhitespace` parameter controls whether content should maintain spacing/formatting
- Array functions typically include overflow handling for large datasets
- Icon functions use the SVG icon system via `components.SVGRef()`
