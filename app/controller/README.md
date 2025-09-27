# Controller Package

This package provides the core HTTP request handling infrastructure for the application.

## Act Function

The `Act` function is the primary entry point for handling HTTP requests. It provides a unified request processing pipeline with telemetry, (optional) authentication, and error handling.

### Signature

```go
func Act(key string, w http.ResponseWriter, r *http.Request, f ActFn)
```

Where `ActFn` is defined as:
```go
type ActFn func(as *app.State, ps *cutil.PageState) (string, error)
```

### Usage Example

```go
func Home(w http.ResponseWriter, r *http.Request) {
    Act("home", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.SetTitleAndData(util.AppName, homeContent)
        return Render(r, as, &views.Home{}, ps)
    })
}
```

## PageState Elements

The `PageState` struct contains numerous properties that influence request processing:

### Authentication & Authorization (when enabled)
- `Admin`: Whether user has admin privileges (bypasses user.Check)
- `Authed`: Authentication status
- `Profile`: User profile information
- `Accounts`: Associated user accounts

### Request Metadata
- `Action`: Current action identifier
- `Method`: HTTP method
- `URI`: Request URL
- `Params`: Filtered request parameters
- `RequestBody`: Raw request body bytes

### Response Control
- `ForceRedirect`: Forces redirect to specified path instead of executing handler
- `Title`: Page title for rendering
- `Description`: Page description
- `Data`: Arbitrary data payload for views

### UI Configuration
- `Menu`: Navigation menu items
- `Breadcrumbs`: Navigation breadcrumbs
- `HideHeader`: Hides page header
- `HideMenu`: Hides navigation menu
- `NoStyle`: Disables standard CSS
- `NoScript`: Disables standard JavaScript
- `Icons`: CSS/UI icons to include
- `RootIcon`, `RootPath`, `RootTitle`: Root navigation configuration

### Browser/Platform Detection
- `Browser`, `BrowserVersion`: Browser identification
- `OS`, `OSVersion`: Operating system detection
- `Platform`: Platform identifier

### Session & State
- `Session`: Session data storage
- `Flashes`: Flash messages for display
- `Context`: Request context for telemetry/cancellation
- `Logger`: Request-scoped logger

### Performance Tracking
- `Started`: Request start time
- `RenderElapsed`: View rendering duration
- `ResponseBytes`: Response size in bytes
