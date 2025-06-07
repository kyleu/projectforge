# OpenAPI

The **`openapi`** module provides Swagger UI integration for [Project Forge](https://projectforge.dev) applications. It embeds the complete Swagger UI interface to visualize, explore, and test OpenAPI specifications.

## Overview

This module provides:

- **Swagger UI Integration**: Complete embedded Swagger UI for interactive API documentation
- **OpenAPI Specification Support**: Compatible with OpenAPI 3.x specifications
- **Interactive Testing**: Built-in API testing capabilities through the Swagger UI
- **Customizable Integration**: Flexible controller implementation with configurable security

## Key Features

### Interactive Documentation
- Visual API documentation with request/response schemas
- Interactive API testing directly in the browser
- Parameter validation and example generation
- Multiple response format support

### Seamless Integration
- Embeds all necessary Swagger UI assets (CSS, JavaScript)
- Works with any OpenAPI 3.x specification
- Supports both local and remote specification URLs
- Integrates with Project Forge's theming system

### Security & Customization
- Flexible security model - add authentication as needed
- Customizable base URL and specification location
- Integration with Project Forge's routing and middleware

## Usage

### Basic Implementation

By default, no controller exposes the Swagger UI. Create a controller following this example, then add a route with your preferred security:

```go
func OpenAPI(w http.ResponseWriter, r *http.Request) {
	controller.Act("openapi", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		u := "http://localhost:{{{ .Port }}}/assets/openapi.json"
		ps.SetTitleAndData("OpenAPI", u)
		return controller.Render(r, as, &vopenapi.OpenAPI{URL: u}, ps, "breadcrumbs")
	})
}
```

### Configuration Options

You can customize the OpenAPI specification source:

```go
// Local specification file
u := "/assets/openapi.json"

// Remote specification
u := "https://api.example.com/openapi.json"

// Dynamic specification
u := fmt.Sprintf("http://%s/api/openapi.json", r.Host)
```

## Configuration

This module requires:

1. An OpenAPI 3.x specification file accessible at the configured URL
2. A controller implementation (not provided by default)
3. Route configuration for accessing the Swagger UI

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/openapi
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [OpenAPI Specification](https://swagger.io/specification/) - OpenAPI 3.x documentation
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/) - Swagger UI features and configuration
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
