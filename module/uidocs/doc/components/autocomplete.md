# Autocomplete

A powerful, server-driven autocomplete component that enhances text input fields with intelligent suggestions. Built on the browser's native `datalist` functionality, this component provides fast, accessible autocomplete features without requiring heavy JavaScript frameworks.

## Key Features

- **Server-Driven**: Fetches suggestions from your backend API in real-time
- **Native Performance**: Uses browser's built-in `datalist` for optimal performance
- **Accessible**: Full keyboard navigation and screen reader support
- **Flexible Data**: Works with any JSON API endpoint
- **Custom Formatting**: Transform server data for display and submission
- **Progressive Enhancement**: Works as regular text input without JavaScript
- **Lightweight**: Minimal JavaScript footprint
- **Search Integration**: Automatically integrates with Project Forge's search system

## How It Works

The autocomplete component follows a simple request-response cycle:

1. **User Types**: As the user types in the input field, the component captures keystrokes
2. **API Request**: After a brief delay, it sends a request to your specified endpoint
3. **Server Response**: Your server returns a JSON array of matching objects
4. **Display Suggestions**: The component formats the data and shows suggestions using a `datalist`
5. **User Selection**: User can select a suggestion or continue typing custom text

## Basic Usage

### Step 1: Create the HTML Input

Start with a standard HTML input element:

```html
<input id="user-search" name="user_id" placeholder="Search for a user..." />
```

### Step 2: Initialize the Autocomplete

Add JavaScript to transform the input into an autocomplete field:

```html
<script>
document.addEventListener("DOMContentLoaded", function() {
    const input = document.getElementById("user-search");

    // Define how to display suggestions
    const title = function(user) {
        return user.name + " (" + user.email + ")";
    };

    // Define what value to submit
    const value = function(user) {
        return user.id;
    };

    // Initialize autocomplete
    // autocomplete(element, url, queryParam, titleFunction, valueFunction)
    {your project}.autocomplete(input, "/api/users/search", "q", title, value);
});
</script>
```

## API Reference

### TypeScript Function Signature

```typescript
function autocomplete(el: HTMLInputElement, url: string, field: string, title: (x: unknown) => string, val: (x: unknown) => string) {
```

**Parameters:**
- `el` (HTMLInputElement): The input element to enhance
- `url` (string): API endpoint URL for fetching suggestions
- `field` (string): Query parameter name for the search term
- `title` (function): Function to format display text from API objects
- `val` (function): Function to extract submission value from API objects

## Server Response Format

Your API endpoint must return a JSON array of objects. The structure of each object is flexible, but should contain the data needed by your title and value functions:

### Simple User Search
```json
[
    {
        "id": "123",
        "name": "John Doe",
        "email": "john@example.com"
    },
    {
        "id": "456",
        "name": "Jane Smith",
        "email": "jane@example.com"
    }
]
```

### Location Search
```json
[
    {
        "id": "city-nyc",
        "name": "New York City",
        "state": "New York",
        "country": "United States",
        "population": 8336817,
        "coordinates": {"lat": 40.7128, "lng": -74.0060}
    }
]
```

### Form Integration

Autocomplete works seamlessly with Project Forge form helpers:

```html
{%- import "views/components" -%}

<form method="post" action="/orders">
    <div class="form-group">
        <label for="customer-search">Customer:</label>
        {%= components.FormInput("customer_id", "customer-search", order.CustomerID, "Search customers...") %}
    </div>

    <div class="form-group">
        <label for="product-search">Product:</label>
        {%= components.FormInput("product_id", "product-search", order.ProductID, "Search products...") %}
    </div>

    <button type="submit">Create Order</button>
</form>

<script>
document.addEventListener("DOMContentLoaded", function() {
    // Customer autocomplete
    const customerInput = document.getElementById("customer-search");
    const customerTitle = function(c) { return c.name + " (" + c.email + ")"; };
    const customerValue = function(c) { return c.id; };
    {your project}.autocomplete(customerInput, "/api/customers/search", "q", customerTitle, customerValue);

    // Product autocomplete
    const productInput = document.getElementById("product-search");
    const productTitle = function(p) { return p.name + " - $" + p.price; };
    const productValue = function(p) { return p.id; };
    {your project}.autocomplete(productInput, "/api/products/search", "q", productTitle, productValue);
});
</script>
```
