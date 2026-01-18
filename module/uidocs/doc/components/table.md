# Tables

Table components that provide styled, interactive, and accessible data tables without requiring JavaScript. These helpers streamline the creation of data tables with features like sorting, resizing, and form integration.

## Key Features

- **No JavaScript Required**: Full functionality using pure CSS and HTML
- **Flexible Styling**: Consistent with Project Forge's design system
- **Performance**: Efficient rendering for large datasets

## Prerequisites

You'll need to import `views/components` at the top of your template to use these table helpers.

```go
{%- import "{your-project}/views/components" -%}
```

## Table Headers

### Basic Table Header

The `TableHeader` component creates interactive column headers with sorting, resizing, and styling capabilities.

```go
{%= components.TableHeader(section, key, title, params, icon, uri, tooltip, sortable, cls, resizable, ps) %}
```

**Parameters:**
- `section`: string - The table section identifier
- `key`: string - The column key for sorting and identification
- `title`: string - Display text for the column header
- `params`: map - URL parameters for sorting and filtering
- `icon`: string - Optional icon to display in the header
- `uri`: string - Base URI for sorting links
- `tooltip`: string - Tooltip text for the header
- `sortable`: bool - Whether the column can be sorted
- `cls`: string - Additional CSS classes
- `resizable`: bool - Whether the column can be resized
- `ps`: PageState - Page state object

**Example:**
```go
<table class="table">
  <thead>
    <tr>
      {%= components.TableHeader("users", "name", "Full Name", params, "", "/users", "Sort by name", true, "", true, ps) %}
      {%= components.TableHeader("users", "email", "Email Address", params, "mail", "/users", "Sort by email", true, "", true, ps) %}
      {%= components.TableHeader("users", "created", "Created", params, "calendar", "/users", "Sort by creation date", true, "text-right", false, ps) %}
      {%= components.TableHeader("users", "actions", "Actions", params, "", "", "Available actions", false, "text-center", false, ps) %}
    </tr>
  </thead>
  <tbody>
    <!-- Table rows go here -->
  </tbody>
</table>
```

## Table Form Inputs

### Basic Table Input

The `TableInput` component creates form inputs within table cells, complete with labels, help text, and proper styling.

```go
{%= components.TableInput(key, id, title, value, indent, helpText) %}
```

**Parameters:**
- `key`: string - Form field name
- `id`: string - DOM element ID
- `title`: string - Label text for the input
- `value`: interface{} - Current value
- `indent`: int - Indentation level for nested inputs
- `helpText`: string - Help text displayed below the input

**Example:**
```go
<table class="form-table">
  <tbody>
    <tr>
      <td>{%= components.TableInput("username", "user-name", "Username", user.Username, 3, "Choose a unique username") %}</td>
    </tr>
    <tr>
      <td>{%= components.TableInput("email", "user-email", "Email", user.Email, 3, "We'll never share your email") %}</td>
    </tr>
  </tbody>
</table>
```

## Advanced Table Features

### Sortable Columns

When `sortable` is set to `true`, the table header becomes clickable and generates appropriate sorting URLs:

```go
// Current URL: /users?sort=name&order=asc
{%= components.TableHeader("users", "email", "Email", params, "", "/users", "Sort by email", true, "", false, ps) %}
// Clicking will navigate to: /users?sort=email&order=asc
```

### Resizable Columns

When `resizable` is set to `true`, users can drag the column borders to adjust width:

```go
{%= components.TableHeader("users", "name", "Name", params, "", "/users", "", true, "", true, ps) %}
```

### Column Icons

Add visual context to columns with icons:

```go
{%= components.TableHeader("users", "email", "Email", params, "mail", "/users", "Email address", true, "", false, ps) %}
{%= components.TableHeader("orders", "total", "Total", params, "dollar", "/orders", "Order total", true, "text-right", false, ps) %}
{%= components.TableHeader("events", "date", "Date", params, "calendar", "/events", "Event date", true, "", false, ps) %}
```

### Custom CSS Classes

Apply custom styling to specific columns:

```go
{%= components.TableHeader("products", "price", "Price", params, "dollar", "/products", "Product price", true, "text-right font-weight-bold", false, ps) %}
```

## Table Input Variations

The `Table.html` file contains helpers for various input types within table cells. Here are some common patterns:

### Text Inputs in Tables
```go
{%= components.TableInput("field_name", "field-id", "Field Label", currentValue, 0, "Help text for this field") %}
```

### Number Inputs in Tables
```go
{%= components.TableInputNumber("quantity", "item-qty", "Quantity", item.Quantity, 0, "Enter the quantity needed") %}
```

### Select Dropdowns in Tables
```go
{%= components.TableSelect("status", "item-status", "Status", item.Status, statusOptions, statusTitles, 0, "Choose the current status") %}
```

### Checkbox Inputs in Tables
```go
{%= components.TableCheckbox("active", "item-active", "Active", item.Active, 0, "Check to activate this item") %}
```

## Responsive Table Patterns

### Horizontal Scrolling
For tables with many columns, wrap in a scrollable container:

```html
<div class="overflow full-width">
  <table class="table expanded">
    <!-- Table content -->
  </table>
</div>
