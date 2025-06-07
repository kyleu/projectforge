# HTTP Archive

The **`har`** module provides comprehensive HTTP Archive (HAR) file parsing and management capabilities for [Project Forge](https://projectforge.dev) applications. HAR files are standard JSON-formatted archives that capture detailed information about HTTP transactions, making them invaluable for debugging, performance analysis, and API documentation.

## Overview

This module enables applications to:

- **Parse HAR Files**: Load and validate HTTP Archive files from various sources
- **Manage Archives**: List, upload, view, and delete HAR files through a web interface
- **Search & Filter**: Find specific entries within HAR files using advanced filtering
- **Export Capabilities**: Convert HAR entries to cURL commands for testing
- **Web Interface**: Full CRUD operations with responsive UI components

## Key Features

### File Management
- Upload HAR files through web interface or file system
- List all available HAR files with metadata
- Delete individual HAR files
- Trim/optimize HAR files by removing unnecessary data

### Data Analysis
- Parse complete HAR file structure (pages, entries, requests, responses)
- Extract timing information for performance analysis
- Access request/response headers, bodies, and metadata
- Support for cookies, cache information, and server details

### Integration
- Search integration for finding specific HTTP transactions
- Content negotiation support (JSON, HTML, CSV export)
- File system abstraction for flexible storage backends
- RESTful API endpoints for programmatic access

## Package Structure

### Core Library

- **`lib/har/service.go`** - Main service for HAR file operations
  - File loading, saving, and deletion
  - Search and filtering capabilities
  - File system integration

- **`lib/har/log.go`** - HAR log structure and metadata
  - Creator and browser information
  - Top-level HAR file representation
  - Web path generation for URLs

- **`lib/har/entry.go`** - Individual HTTP transaction entries
  - Request/response pair representation
  - Timing and performance data
  - Connection and server information

- **`lib/har/request.go`** - HTTP request details
  - Headers, query parameters, and POST data
  - Method, URL, and HTTP version
  - Cookie and authentication information

- **`lib/har/response.go`** - HTTP response details
  - Status codes and reason phrases
  - Response headers and content
  - Compression and encoding information

### Utilities

- **`lib/har/page.go`** - Page timing and navigation data
- **`lib/har/cookies.go`** - Cookie parsing and management
- **`lib/har/cache.go`** - HTTP cache information
- **`lib/har/curl.go`** - cURL command generation from HAR entries
- **`lib/har/util.go`** - Helper functions and constants
- **`lib/har/selector.go`** - Entry filtering and selection

### Web Interface

- **`controller/clib/har.go`** - HTTP handlers for HAR operations
  - List, upload, view, and delete operations
  - Content negotiation for multiple formats
  - Error handling and validation

- **`controller/routes/har.go`** - Route definitions
  - RESTful endpoint mapping
  - Parameter extraction and validation

- **`controller/cmenu/har.go`** - Navigation menu integration

## Usage Examples

### Loading HAR Files

```go
// Create service with file system backend
service := har.NewService(fs)

// Load a HAR file
log, err := service.Load("example.har")
if err != nil {
    return err
}

// Access entries
for _, entry := range log.Entries {
    fmt.Printf("Request: %s %s\n", entry.Request.Method, entry.Request.URL)
    fmt.Printf("Response: %d\n", entry.Response.Status)
}
```

### Web API Endpoints

- `GET /har` - List all HAR files
- `POST /har` - Upload new HAR file
- `GET /har/{key}` - View specific HAR file details
- `GET /har/{key}/delete` - Delete HAR file
- `GET /har/{key}/trim` - Optimize HAR file size

### Search Integration

The module integrates with Project Forge's search system to find HAR files and entries matching specific criteria.

## Configuration

No additional configuration required beyond standard Project Forge setup. HAR files are stored in the `./har` directory relative to the application root.

## File Format Support

Supports the full [HAR 1.2 specification](http://www.softwareishard.com/blog/har-12-spec/):

- Log metadata (creator, browser, version)
- Page information and timing
- Request/response pairs with full headers
- Cookie and cache information
- Content compression and encoding
- Custom fields and comments

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/har
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Project Forge Documentation](https://projectforge.dev) - Complete documentation  
- [HAR 1.2 Specification](http://www.softwareishard.com/blog/har-12-spec/) - Official HAR format documentation
