# Flash Notifications

Flash notifications provide temporary feedback messages to users after actions like form submissions, data updates, or system events. These messages appear briefly on page load and automatically disappear, providing a non-intrusive way to communicate with users.

## Overview

Flash notifications are essential for user experience, providing immediate feedback about the success or failure of user actions. They appear as styled panels in the upper-right corner of the page and can be dismissed manually or automatically after a timeout.

## Key Features

- **No JavaScript Required**: Full functionality using pure CSS and HTML
- **Temporary Display**: Messages appear on page load and clear before the next request
- **Auto-dismiss**: Messages fade away automatically after a few seconds (requires JavaScript)
- **Manual Dismiss**: Users can close messages immediately with a close button (no JavaScript required)
- **Multiple Types**: Support for success, error, warning, and info message types
- **Non-intrusive**: Positioned to not interfere with page content

## How Flash Messages Work

Flash messages follow the "flash and redirect" pattern common in web applications:

1. User performs an action (submits form, deletes item, etc.)
2. Server processes the action
3. Server sets a flash message indicating success or failure (stored in the response's cookie)
4. Server redirects to prevent duplicate submissions
5. New page loads and displays the flash message
6. Message is cleared from the session after display

## Basic Usage

### Using FlashAndRedir Helper

The most common way to set flash messages is using the `FlashAndRedir` helper function:

```go
// FlashAndRedir(success bool, msg string, redir string, w http.ResponseWriter, ps *cutil.PageState)
return controller.FlashAndRedir(true, "User created successfully!", "/users", w, ps)
```

**Parameters:**
- `success`: bool - Whether this is a success (true) or error (false) message
- `msg`: string - The message text to display
- `redir`: string - URL to redirect to after setting the flash
- `w`: http.ResponseWriter - The HTTP response writer
- `ps`: *cutil.PageState - The current page state

### Manual Flash Messages

You can also add flash messages manually by appending to the `Flashes` field:

```go
// Add a success message
ps.Flashes = append(ps.Flashes, &cutil.Flash{
    Type:    "success",
    Message: "Operation completed successfully!",
})

// Add an error message
ps.Flashes = append(ps.Flashes, &cutil.Flash{
    Type:    "error",
    Message: "Something went wrong. Please try again.",
})

// Add a warning message
ps.Flashes = append(ps.Flashes, &cutil.Flash{
    Type:    "warning",
    Message: "This action cannot be undone.",
})

// Add an info message
ps.Flashes = append(ps.Flashes, &cutil.Flash{
    Type:    "info",
    Message: "New features are now available!",
})
```
