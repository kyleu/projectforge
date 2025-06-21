# Menu

A flexible navigation menu system that provides hierarchical navigation without requiring JavaScript. The menu system supports multi-level navigation, active state management, and responsive design while maintaining full functionality in JavaScript-disabled environments.

## Overview

The menu component is the primary navigation system for Project Forge applications. It provides a clean, accessible way to organize navigation links with support for nested menus, active state highlighting, and responsive behavior across different screen sizes.

## Key Features

- **No JavaScript Required**: Full functionality using pure CSS and HTML
- **Progressive Enhancement**: Enhanced with JavaScript when available for improved UX
- **Hierarchical Navigation**: Support for multi-level menu structures
- **Active State Management**: Automatic highlighting of current page/section
- **Responsive Design**: Adapts to mobile and desktop layouts
- **Keyboard Accessible**: Full keyboard navigation support
- **Auto-scroll**: Active items scroll into view on page load (with JavaScript)
- **Flexible Structure**: Configurable menu items and organization

## How It Works

The menu system uses Go structs to define menu structure and generates HTML with proper CSS classes for styling and behavior. Active states are determined by comparing the current URL path with menu item paths.

## Menu Configuration

Menu structure is defined in Go code, typically in your application's menu configuration:

```go
// Example menu structure
type MenuItem struct {
    Key         string
    Title       string
    Description string
    Icon        string
    Path        string
    Children    []*MenuItem
}

// Define your application menu
var AppMenu = &MenuItem{
    Key:   "root",
    Title: "Main Navigation",
    Children: []*MenuItem{
        {
            Key:         "dashboard",
            Title:       "Dashboard",
            Description: "Application overview and statistics",
            Icon:        "dashboard",
            Path:        "/dashboard",
        },
        {
            Key:         "users",
            Title:       "Users",
            Description: "User management and administration",
            Icon:        "users",
            Path:        "/users",
            Children: []*MenuItem{
                {
                    Key:         "users-list",
                    Title:       "All Users",
                    Description: "View and manage all users",
                    Icon:        "list",
                    Path:        "/users",
                },
                {
                    Key:         "users-create",
                    Title:       "Add User",
                    Description: "Create a new user account",
                    Icon:        "plus",
                    Path:        "/users/new",
                },
            },
        },
        {
            Key:         "settings",
            Title:       "Settings",
            Description: "Application configuration",
            Icon:        "cog",
            Path:        "/settings",
        },
    },
}
```

## Active State Management

The menu system automatically determines active states by comparing the current URL path with menu item paths:

### Current Path Matching
- **Exact Match**: `/users/new` matches menu item with path `/users/new`
- **Prefix Match**: `/users/123` matches menu item with path `/users`
- **Section Match**: Any path starting with `/admin` matches admin section
