# Markdown

Project Forge provides comprehensive utilities for rendering Markdown files as HTML, making it easy to create documentation systems, content management features, and rich text displays. The markdown system is designed to integrate seamlessly with Project Forge's templating and component systems.

## Overview

The markdown utilities offer:
- **No JavaScript Required**: Full functionality using pure CSS and HTML
- **Server-side Rendering**: Convert Markdown to HTML on the server
- **Content Processing**: Clean and format HTML output for safe display
- **Documentation Integration**: Built-in support for documentation browsing
- **Flexible Transformation**: Custom processing pipelines for specialized content

## Basic Usage

```go
import (
    "github.com/gomarkdown/markdown"
    "{your_app}/app/controller/cutil"
)

// Basic markdown to HTML conversion
markdownContent := `
# Welcome to Project Forge

This is a **bold** statement with *italic* text.

## Features

- Fast rendering
- Clean HTML output
- Security-focused processing
`

html := string(markdown.ToHTML([]byte(markdownContent), nil, nil))

// Clean and format the HTML for template use
title, cleanHTML, err := cutil.FormatCleanMarkup(html, "document")
if err != nil {
    // Handle error
}
```

## Content Security and Sanitization

### HTML Sanitization

The `FormatCleanMarkup` function provides built-in HTML sanitization:

```go
// FormatCleanMarkup performs several security and formatting operations:
// 1. Sanitizes HTML to prevent XSS attacks
// 2. Adds appropriate CSS classes
// 3. Processes links and images
// 4. Extracts title from content
// 5. Adds optional icon integration

title, cleanHTML, err := cutil.FormatCleanMarkup(rawHTML, "optional-icon")
```
