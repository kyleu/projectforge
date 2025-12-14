# Search

The **`search`** module provides comprehensive search functionality for your application.
It adds a powerful, extensible search system with a clean UI integrated into the application navigation.

## Overview

This module provides:

- **Global Search Bar**: Integrated search interface in the navigation header
- **Provider-Based Architecture**: Extensible system for searching different data sources
- **Advanced Query Parsing**: Support for general queries and keyed search syntax (`key:value`)
- **Results Highlighting**: Automatic match highlighting with query term detection
- **Smart Navigation**: Auto-redirect for single results, grouped display for multiple matches

## Key Features

### Search Architecture
- Provider-based search system for easy extensibility
- Async result collection from multiple data sources
- Type-safe result structures with metadata
- Configurable result limits and scoring

### User Experience
- Clean, responsive search interface
- Real-time search suggestions
- Match highlighting with context
- Error handling and user feedback
- Mobile-friendly design

### Query Support
- General text search across all providers
- Keyed search syntax: `user:john`, `status:active`
- Multiple query terms with AND logic
- Case-insensitive matching

## Usage

### Basic Implementation

1. **Enable the module** in your application's configuration

2. **Implement search providers** by modifying `./app/lib/search/search.go`:

```go
func Search(ctx context.Context, params *Params, ps *cutil.PageState, logger util.Logger) *Results {
    ret := NewResults(params)

    // Add your search providers
    if params.Q != "" {
        // Search users
        if userResults := searchUsers(ctx, params.Q); len(userResults) > 0 {
            ret.Results = append(ret.Results, userResults...)
        }

        // Search content
        if contentResults := searchContent(ctx, params.Q); len(contentResults) > 0 {
            ret.Results = append(ret.Results, contentResults...)
        }
    }

    return ret
}
```

3. **Create search providers** for your data sources:

```go
func searchUsers(ctx context.Context, query string) result.Results {
    // Implement user search logic
    // Return result.Result instances with proper highlighting
}
```

### Advanced Features

Create typed results for better UX:

```go
func NewUserResult(user *User, query string) *result.Result {
    return &result.Result{
        Type:    "user",
        Title:   user.Name,
        Summary: user.Email,
        URL:     fmt.Sprintf("/users/%d", user.ID),
        Matches: result.HighlightMatches(user.Name, query),
    }
}
```

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/search
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
