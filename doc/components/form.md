# Form Helpers

A comprehensive collection of HTML input template helpers that streamline form creation in Project Forge applications. These components provide consistent styling, validation support, and accessibility features across all form elements.

## Key Features

- **No JavaScript Required**: Full functionality using pure CSS and HTML
- **Consistent Styling**: All form elements follow the same design patterns
- **Accessibility Built-in**: Proper labeling, ARIA attributes, and keyboard navigation
- **Type Safety**: Go template helpers with proper type handling
- **Validation Ready**: Integration with server-side validation systems
- **Multiple Layouts**: Support for different form layouts (horizontal, vertical, table-based)
- **Flexible Options**: Customizable placeholders, validation, and styling

## Prerequisites

You'll need to import `views/components` at the top of your template to use these form helpers.

## Common Parameters

Most form helpers share these common parameters:

```go
// name of the form field (used for form submission)
var key = "username"

// optional DOM element id (for CSS targeting and accessibility)
var id = "input-username"

// current value for the input (type varies by input type)
var value = user.Username

// placeholder text shown when field is empty
var placeholder = "Enter your username"
```

## Text Input Components

### Basic Text Input
For general text input fields like names, titles, or descriptions.

```go
{%= components.FormInput(key, id, value, placeholder) %}
```

**Example:**
```go
{%= components.FormInput("full_name", "user-name", user.FullName, "Enter your full name") %}
```

### Password Input
For password fields that hide the input text.

```go
{%= components.FormInputPassword(key, id, value, placeholder) %}
```

**Example:**
```go
{%= components.FormInputPassword("password", "user-password", "", "Enter your password") %}
```

### Multi-line Text Area
For longer text content like descriptions, comments, or messages.

```go
{%= components.FormTextarea(key, id, value, placeholder) %}
```

**Example:**
```go
{%= components.FormTextarea("description", "product-desc", product.Description, "Describe your product...") %}
```

## Numeric Input Components

### Integer Numbers
For whole number inputs like quantities, ages, or counts.

```go
{%= components.FormInputNumber(key, id, value, placeholder) %}
```

**Example:**
```go
{%= components.FormInputNumber("quantity", "item-qty", order.Quantity, "Enter quantity") %}
```

### Floating Point Numbers
For decimal numbers like prices, measurements, or percentages.

```go
{%= components.FormInputFloat(key, id, value, placeholder) %}
```

**Example:**
```go
{%= components.FormInputFloat("price", "product-price", product.Price, "0.00") %}
```

## Date and Time Components

### Date and Time Input
For full timestamp input including both date and time components.

```go
{%= components.FormInputTimestamp(key, id, value, placeholder) %}
```

**Example:**
```go
{%= components.FormInputTimestamp("created_at", "event-time", event.CreatedAt, "") %}
```

### Date Only Input
For date-only input without time component.

```go
{%= components.FormInputTimestampDay(key, id, value, placeholder) %}
```

**Example:**
```go
{%= components.FormInputTimestampDay("birth_date", "user-birthday", user.BirthDate, "") %}
```

## Specialized Input Components

### UUID Input
For UUID fields with proper validation and formatting.

```go
{%= components.FormInputUUID(key, id, value, placeholder) %}
```

**Example:**
```go
{%= components.FormInputUUID("user_id", "related-user", relation.UserID, "") %}
```

## Selection Components

### Dropdown Select
For single-choice selection from a predefined list of options.

```go
{%= components.FormSelect(key, id, value, opts, titles, placeholder) %}
```

**Parameters:**
- `opts`: []string - The option values
- `titles`: []string - Display text for each option (can be same as opts)

**Example:**
```go
var statusOpts = []string{"active", "inactive", "pending"}
var statusTitles = []string{"Active", "Inactive", "Pending"}
{%= components.FormSelect("status", "user-status", user.Status, statusOpts, statusTitles, "Select status") %}
```

### Autocomplete Input (Datalist)
For text input with suggested options that allows free-form entry.

```go
{%= components.FormDatalist(key, id, value, opts, titles, placeholder) %}
```

**Example:**
```go
var cityOpts = []string{"new-york", "los-angeles", "chicago"}
var cityTitles = []string{"New York", "Los Angeles", "Chicago"}
{%= components.FormDatalist("city", "user-city", user.City, cityOpts, cityTitles, "Enter or select city") %}
```

### Radio Buttons
For single-choice selection with all options visible.

```go
{%= components.FormRadio(key, value, opts, titles, indent) %}
```

**Parameters:**
- `indent`: int - Indentation level for styling

**Example:**
```go
var themeOpts = []string{"light", "dark", "auto"}
var themeTitles = []string{"Light Mode", "Dark Mode", "Auto"}
{%= components.FormRadio("theme", user.Theme, themeOpts, themeTitles, 0) %}
```

### Checkboxes
For multiple-choice selection where users can select several options.

```go
{%= components.FormCheckbox(key, values, opts, titles, indent) %}
```

**Parameters:**
- `values`: []string - Currently selected values
- `indent`: int - Indentation level for styling

**Example:**
```go
var interestOpts = []string{"tech", "sports", "music", "travel"}
var interestTitles = []string{"Technology", "Sports", "Music", "Travel"}
{%= components.FormCheckbox("interests", user.Interests, interestOpts, interestTitles, 0) %}
```

## Layout Variations

### Standard Form Layout
The default form helpers create horizontally-aligned forms suitable for most use cases.

### Vertical Form Layout
For forms where labels should appear above inputs:

```go
{%= components.FormVerticalInput(key, id, value, placeholder) %}
```

### Table-based Form Layout
For forms that need to be embedded within tables or need compact spacing:

```go
{%= components.TableInput(key, id, title, value, indent, "help text") %}
```

**Example:**
```go
{%= components.TableInput("email", "user-email", "Email Address", user.Email, 0, "We'll never share your email") %}
```

## Complete Form Example

Here's a comprehensive example showing various form components working together:

```html
<form action="/users" method="post" class="form">
  <fieldset>
    <legend>User Information</legend>

    <!-- Basic text inputs -->
    <div class="form-group">
      <label for="user-name">Full Name:</label>
      {%= components.FormInput("full_name", "user-name", user.FullName, "Enter your full name") %}
    </div>

    <div class="form-group">
      <label for="user-email">Email:</label>
      {%= components.FormInput("email", "user-email", user.Email, "your@email.com") %}
    </div>

    <!-- Password input -->
    <div class="form-group">
      <label for="user-password">Password:</label>
      {%= components.FormInputPassword("password", "user-password", "", "Enter a secure password") %}
    </div>

    <!-- Date input -->
    <div class="form-group">
      <label for="user-birthday">Birth Date:</label>
      {%= components.FormInputTimestampDay("birth_date", "user-birthday", user.BirthDate, "") %}
    </div>

    <!-- Select dropdown -->
    <div class="form-group">
      <label for="user-role">Role:</label>
      {%= components.FormSelect("role", "user-role", user.Role,
          []string{"user", "admin", "moderator"},
          []string{"User", "Administrator", "Moderator"},
          "Select a role") %}
    </div>

    <!-- Radio buttons -->
    <fieldset class="form-group">
      <legend>Theme Preference:</legend>
      {%= components.FormRadio("theme", user.Theme,
          []string{"light", "dark", "auto"},
          []string{"Light", "Dark", "Auto"}, 1) %}
    </fieldset>

    <!-- Checkboxes -->
    <fieldset class="form-group">
      <legend>Interests:</legend>
      {%= components.FormCheckbox("interests", user.Interests,
          []string{"tech", "sports", "music"},
          []string{"Technology", "Sports", "Music"}, 1) %}
    </fieldset>

    <!-- Textarea -->
    <div class="form-group">
      <label for="user-bio">Biography:</label>
      {%= components.FormTextarea("bio", "user-bio", user.Bio, "Tell us about yourself...") %}
    </div>

    <!-- Submit button -->
    <div class="form-actions">
      <button type="submit" class="btn btn-primary">Save User</button>
      <a href="/users" class="btn btn-secondary">Cancel</a>
    </div>
  </fieldset>
</form>
```
