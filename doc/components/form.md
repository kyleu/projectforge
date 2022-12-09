# Form Helpers

HTML input template helpers. You'll need to import `views/components` at the top of your template.

The component that renders a normal HTML input field is named `FormInput`. 
For each component, alternate version or available, such as `TableInput` for table-based form fields, and `FormVerticalInput` for a vertical stacked form.

```go
// name of the form field
var key = "name"

// optional DOM element id
var id = "input-name"

// value for the input, type varies
var value = obj.Field

// only used for some inputs, shows when value is empty
var placeholder = "Type Stuff" 
```

Basic text input, string value
```go
{%= components.FormInput(key, id, value, placeholder) %}
```

Password input, string value (hidden in UI)
```go
{%= components.FormInputPassword(key, id, value, placeholder) %}
```

Integer number input, value can be any scalar type
```go
{%= components.FormInputNumber(key, id, value, placeholder) %}
```

Integer number input, value can be any scalar type
```go
{%= components.FormInputFloat(key, id, value, placeholder) %}
```

Date and time input, value is *time.Time
```go
{%= components.FormInputTimestamp(key, id, value, placeholder) %}
```

Date input (no time component), value is *time.Time
```go
{%= components.FormInputTimestampDay(key, id, value, placeholder) %}
```

UUID input, *uuid.UUID value
```go
{%= components.FormInputUUID(key, id, value, placeholder) %}
```

Multi-line text input via HTML textarea, string value
```go
{%= components.FormTextarea(key, id, value, placeholder) %}
```

X input, string value
```go
{%= components.FormInputX(key, id, value, placeholder) %}
```

Single-choice select, string value, accepts []string for options and titles
```go
{%= components.FormSelect(key, id, value, opts, titles, placeholder) %}
```

Autocompleted text input, string value, accepts []string for options and titles, allows all inputs
```go
{%= components.FormDatalist(key, id, value, opts, titles, placeholder) %}
```

Single-choice radio inputs, string value, accepts []string for options and titles
```go
{%= components.FormRadio(key, value, opts, titles, indent) %}
```

Multiple-choice checkbox inputs, []string values, accepts []string for options and titles
```go
{%= components.FormCheckbox(key, values, opts, titles, indent) %}
```
