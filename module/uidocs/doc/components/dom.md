# DOM Utilities

Project Forge provides a comprehensive set of DOM utilities that make working with HTML elements more convenient and type-safe. These utilities are designed to work without requiring external dependencies and provide a clean, functional approach to DOM manipulation.

## Overview

The DOM utilities offer three main categories of functionality:
- **Element Selection**: Find elements with type safety and error handling
- **Content Manipulation**: Safely update element content and properties
- **Display Control**: Show, hide, and manage element visibility

## Element Selection Functions

### `els<T>(selector, context?)` - Find Multiple Elements

Returns an array of all elements matching the CSS selector. This function is type-safe and will return an empty array if no elements are found.

```typescript
import {els} from "./dom";

// Find all buttons on the page
const allButtons = els<HTMLButtonElement>("button");

// Find all buttons within a specific container
const containerButtons = els<HTMLButtonElement>("button", someContainer);

// Find elements by class
const cards = els<HTMLDivElement>(".card");

// Find elements by attribute
const requiredInputs = els<HTMLInputElement>("input[required]");
```

**Parameters:**
- `selector`: CSS selector string
- `context` (optional): Parent element to search within

**Returns:** Read-only array of typed HTML elements

### `opt<T>(selector, context?)` - Find Single Element (Optional)

Returns a single element or `undefined` if not found. This is perfect when you're not sure if an element exists and want to handle the absence gracefully.

```typescript
import {opt} from "./dom";

// Safely get an element that might not exist
const sidebar = opt<HTMLDivElement>("#sidebar");
if (sidebar) {
  // Element exists, safe to use
  sidebar.classList.add("active");
}

// Check for optional form elements
const optionalField = opt<HTMLInputElement>("#optional-email");
if (optionalField && optionalField.value) {
  // Process the optional field only if it exists and has a value
  processEmail(optionalField.value);
}
```

**Parameters:**
- `selector`: CSS selector string
- `context` (optional): Parent element to search within

**Returns:** Typed HTML element or `undefined`

**Note:** If multiple elements match the selector, a warning will be logged to the console.

### `req<T>(selector, context?)` - Find Single Element (Required)

Returns a single element or throws an error if not found. Use this when you expect the element to exist and want to fail fast if it doesn't.

```typescript
import {req} from "./dom";

// Get an element that must exist
const mainContent = req<HTMLDivElement>("#main-content");

// Get form elements that are required for functionality
const loginForm = req<HTMLFormElement>("#login-form");
const submitButton = req<HTMLButtonElement>("#submit-btn");

// Will throw an error if element doesn't exist
try {
  const criticalElement = req<HTMLDivElement>("#critical-component");
  initializeCriticalFeature(criticalElement);
} catch (error) {
  console.error("Critical element missing:", error.message);
}
```

**Parameters:**
- `selector`: CSS selector string
- `context` (optional): Parent element to search within

**Returns:** Typed HTML element

**Throws:** Error if no element is found

## Content Manipulation Functions

### `setHTML(element, html)` - Set Inner HTML

Safely sets the innerHTML of an element. Accepts either an element reference or a CSS selector string.

```typescript
import {setHTML} from "./dom";

// Using element reference
const container = req<HTMLDivElement>("#container");
setHTML(container, "<p>New content with <strong>HTML</strong></p>");

// Using selector string
setHTML("#status", "<span class='success'>✓ Complete</span>");

// Dynamic content generation
const items = ["Apple", "Banana", "Cherry"];
const listHTML = items.map(item => `<li>${item}</li>`).join("");
setHTML("#fruit-list", listHTML);
```

### `setText(element, text)` - Set Text Content

Sets the text content of an element, automatically escaping HTML characters for security.

```typescript
import {setText} from "./dom";

// Safe text setting (HTML will be escaped)
setText("#user-name", userInput); // Safe even if userInput contains HTML

// Update status messages
setText("#status-message", "Processing complete!");

// Using with element reference
const title = req<HTMLHeadingElement>("h1");
setText(title, "Welcome to Project Forge");
```

### `clear(element)` - Clear Element Content

Removes all content from an element by setting its innerHTML to an empty string.

```typescript
import {clear} from "./dom";

// Clear a container before adding new content
clear("#results");

// Clear form validation messages
clear(".error-messages");

// Using with element reference
const notifications = req<HTMLDivElement>("#notifications");
clear(notifications);
```

## Display Control Functions

### `setDisplay(element, condition, displayValue?)` - Control Element Visibility

Shows or hides elements based on a boolean condition. This is more convenient than manually setting CSS display properties.

```typescript
import {setDisplay} from "./dom";

// Basic show/hide
setDisplay("#loading-spinner", isLoading);
setDisplay("#error-message", hasError);

// Custom display value (default is "block")
setDisplay("#sidebar", isOpen, "flex");
setDisplay("#inline-element", shouldShow, "inline-block");

// Conditional UI updates
const user = getCurrentUser();
setDisplay("#login-section", !user);
setDisplay("#user-profile", !!user);
setDisplay("#admin-panel", user?.isAdmin);
```

**Parameters:**
- `element`: Element reference or CSS selector string
- `condition`: Boolean determining visibility
- `displayValue` (optional): CSS display value when shown (default: "block")

## Usage Patterns

### Form Handling

```typescript
import {req, opt, setText, setDisplay} from "./dom";

function validateForm() {
  const form = req<HTMLFormElement>("#contact-form");
  const emailInput = req<HTMLInputElement>("#email");
  const errorDiv = opt<HTMLDivElement>("#email-error");

  const isValid = emailInput.value.includes("@");

  if (errorDiv) {
    setText(errorDiv, isValid ? "" : "Please enter a valid email");
    setDisplay(errorDiv, !isValid);
  }

  return isValid;
}
```

### Dynamic Content Loading

```typescript
import {req, setHTML, setDisplay} from "./dom";

async function loadUserProfile(userId: string) {
  const container = req<HTMLDivElement>("#profile-container");
  const loader = req<HTMLDivElement>("#profile-loader");

  setDisplay(loader, true);

  try {
    const userData = await fetchUser(userId);
    const profileHTML = generateProfileHTML(userData);
    setHTML(container, profileHTML);
  } catch (error) {
    setHTML(container, "<p class='error'>Failed to load profile</p>");
  } finally {
    setDisplay(loader, false);
  }
}
```
