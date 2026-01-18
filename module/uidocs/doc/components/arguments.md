# Arguments

A powerful component system for collecting and validating user input before executing actions in Project Forge applications. The Arguments component provides a clean, reusable way to gather required parameters from users through web forms with built-in validation and error handling.

## Key Features

- **Declarative Field Definition**: Define fields with types, validation, and defaults
- **Automatic Form Generation**: Creates properly styled forms from field definitions
- **Built-in Validation**: Server-side validation with user-friendly error messages
- **Type Safety**: Proper type handling for different input types (text, numbers, dates, etc.)
- **Progressive Enhancement**: Works without JavaScript, enhanced with it
- **Consistent Styling**: Integrates seamlessly with Project Forge's design system
- **Flexible Workflow**: Collect arguments first, then execute actions with validated data

## How It Works

The Arguments component follows a simple workflow:

1. **Define Fields**: Create a `util.FieldDescs` structure describing your required inputs
2. **Collect Input**: Use `util.FieldDescsCollect()` to gather and validate user input
3. **Check Completeness**: Test if all required fields are present and valid
4. **Show Form or Execute**: Display the form for missing fields, or proceed with the action

## Basic Usage

### Step 1: Define Your Fields

Create a `util.FieldDescs` structure that describes the arguments you need:

```go
package controllers

var orderArgs = util.FieldDescs{
	{Key: "name", Title: "Customer Name", Description: "Enter the customer's full name"},
	{Key: "quantity", Title: "Quantity", Description: "Number of items to order", Type: "number", Default: "1"},
	{Key: "priority", Title: "Priority Level", Description: "Order priority", Type: "select", Options: []string{"low", "normal", "high"}, Default: "normal"},
}
```

### Step 2: Implement Your Controller

Use the arguments in your HTTP handler:

```go
func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	controller.Act("place.order", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		// Collect and validate arguments from the request
		argRes := util.FieldDescsCollect(r, orderArgs)

		// Check if any required fields are missing or invalid
		if argRes.HasMissing() {
			// Set up the page to display the argument collection form
			ps.SetTitleAndData("Order Options", argRes)
			msg := "Please provide the following information to place your order"
			return controller.Render(r, as, &vpage.Args{
				URL: r.URL.String(),
				Directions: msg,
				ArgRes: argRes
			}, ps, "breadcrumb")
		}

		// All arguments are present and valid - proceed with the action
		ord := actuallyPlaceOrder(argRes.Values)
		return "/order/" + ord.ID.String(), nil
	})
}
```

## Field Definition Reference

### Basic Field Properties

```go
util.FieldDesc{
	Key:         "field_name",        // Form field name (required)
	Title:       "Display Name",      // Human-readable label (required)
	Description: "Help text",         // Additional guidance for users
	Type:        "text",              // Input type (see types below)
	Default:     "default_value",     // Pre-filled value
	Required:    true,                // Whether field is mandatory
}
```

### Supported Field Types

#### Text Input Types
```go
// Basic text input
{Key: "name", Title: "Full Name", Type: "text"}

// Email input with validation
{Key: "email", Title: "Email Address", Type: "email"}

// Password input (hidden text)
{Key: "password", Title: "Password", Type: "password"}

// Multi-line text area
{Key: "description", Title: "Description", Type: "textarea"}

// URL input with validation
{Key: "website", Title: "Website URL", Type: "url"}
```

#### Numeric Input Types
```go
// Integer numbers
{Key: "age", Title: "Age", Type: "number", Default: "18"}

// Floating point numbers
{Key: "price", Title: "Price", Type: "float", Default: "0.00"}

// Numbers with min/max constraints
{Key: "quantity", Title: "Quantity", Type: "number", Min: "1", Max: "100"}
```

#### Date and Time Types
```go
// Date only
{Key: "birth_date", Title: "Birth Date", Type: "date"}

// Date and time
{Key: "appointment", Title: "Appointment Time", Type: "datetime"}

// Time only
{Key: "meeting_time", Title: "Meeting Time", Type: "time"}
```

#### Selection Types
```go
// Dropdown select
{Key: "status", Title: "Status", Type: "select",
	Options: []string{"active", "inactive", "pending"},
	Titles: []string{"Active", "Inactive", "Pending"}}

// Radio buttons
{Key: "priority", Title: "Priority", Type: "radio",
	Options: []string{"low", "normal", "high"},
	Titles: []string{"Low Priority", "Normal Priority", "High Priority"}}

// Checkboxes (multiple selection)
{Key: "features", Title: "Features", Type: "checkbox",
	Options: []string{"ssl", "backup", "monitoring"},
	Titles: []string{"SSL Certificate", "Daily Backup", "24/7 Monitoring"}}
```

#### Boolean Type
```go
// Single checkbox for yes/no
{Key: "agree_terms", Title: "I agree to the terms", Type: "bool"}
```

## Working with Collected Arguments

### Accessing Values
```go
argRes := util.FieldDescsCollect(r, myArgs)

// Get a specific value
customerName := argRes.Values["customer_name"]

// Get with type conversion
quantity, err := argRes.GetInt("quantity")
price, err := argRes.GetFloat("price")
isRushOrder, err := argRes.GetBool("rush_order")
deliveryDate, err := argRes.GetTime("delivery_date")

// Get with defaults
quantity := argRes.GetIntDefault("quantity", 1)
```

### Validation and Error Handling
```go
argRes := util.FieldDescsCollect(r, myArgs)

// Check for missing required fields
if argRes.HasMissing() {
	// Show the form with error messages
	return showArgumentForm(argRes)
}

// Check for validation errors
if argRes.HasErrors() {
	// Handle specific validation errors
	for field, error := range argRes.Errors {
		log.Printf("Field %s has error: %s", field, error)
	}
	return showArgumentForm(argRes)
}

// All good - proceed with action
return executeAction(argRes.Values)
```
