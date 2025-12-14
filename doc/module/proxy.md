# Proxy

The **`proxy`** module provides secure HTTP proxy functionality for your application.
It allows your application to proxy requests to external services while maintaining security controls and consistent routing patterns.

## Overview

This module enables applications to act as a secure HTTP proxy, forwarding requests to configured external services while:

- **Enforcing Security**: Maintains your application's authentication and authorization
- **URL Rewriting**: Automatically rewrites relative URLs in responses to maintain proper routing
- **Header Management**: Filters and manages HTTP headers for security
- **Service Management**: Dynamically register and manage proxy destinations

⚠️ **Security Notice**: This module is marked as "dangerous" as it can expose your application to external services. Use with proper authentication and validation.

## Key Features

### Secure Proxying
- Filters dangerous HTTP headers (Connection, Keep-Alive, etc.)
- Adds "Proxied" header to identify proxy requests
- Maintains request context and logging

### URL Rewriting
- Automatically rewrites `href="/..."` and `src="/..."` attributes
- Ensures proxied content links work correctly within your application
- Maintains relative path integrity

### Dynamic Service Management
- Register and remove proxy destinations at runtime
- List all configured proxy services
- Per-service URL configuration

## Configuration

The proxy service is initialized with:

```go
// Create proxy service with URL prefix and initial services
proxyService := proxy.NewService("/proxy", map[string]string{
    "api": "http://localhost:8080",
    "docs": "https://docs.example.com",
})

// Add to application state
as.Services.Proxy = proxyService
```

## Usage

### Basic Setup

1. **Initialize the service** in your application startup:
```go
proxy := proxy.NewService("/proxy", initialProxies)
```

2. **Wire up routes** in your router configuration:
```go
// Add proxy routes (typically in routes/proxy.go)
mux.HandleFunc("/proxy", clib.ProxyIndex)
mux.HandleFunc("/proxy/{svc}/{path:.*}", clib.ProxyHandle)
```

### Managing Proxy Services

```go
// Add a new service
proxyService.SetURL("newservice", "https://api.example.com")

// Remove a service
proxyService.Remove("oldservice")

// List all services
services := proxyService.List()
```

### Request Flow

1. Client makes request to `/proxy/{service}/{path}`
2. Proxy service looks up the target URL for `{service}`
3. Request is forwarded to `{target_url}/{path}`
4. Response is processed and URLs are rewritten
5. Modified response is returned to client

## Security Considerations

- **Authentication**: Proxy requests should be authenticated through your application's normal auth flow
- **Authorization**: Consider implementing service-specific access controls
- **Input Validation**: Validate service names and paths to prevent abuse
- **Network Access**: Ensure proxied services are on trusted networks
- **Header Filtering**: The module filters dangerous headers but review for your specific needs

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/proxy
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [HTTP Security Best Practices](https://developer.mozilla.org/en-US/docs/Web/HTTP/Security)
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
- [Configuration Variables](../running.md) - Available environment variables
