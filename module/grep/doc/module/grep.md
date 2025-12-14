# Grep Search

The **`grep`** module provides fast text search capabilities for your application using the high-performance [ripgrep](https://github.com/BurntSushi/ripgrep) command-line tool.

## Overview

This module enables applications to perform blazing-fast text searches across files and directories by wrapping the `ripgrep` utility in a convenient Go API. It's particularly useful for:

- **Code search**: Finding functions, variables, or patterns across codebases
- **Log analysis**: Searching through application logs and system files
- **Content discovery**: Locating specific text within large document collections
- **Development tools**: Building IDE-like search functionality

## Key Features

### Performance
- Leverages ripgrep's Rust-based regex engine for maximum speed
- Handles large codebases and file collections efficiently
- Supports recursive directory traversal with smart filtering

### Search Capabilities
- Regular expression support with multiple regex engines
- Case-sensitive and case-insensitive searching
- File type filtering and glob pattern matching
- Configurable context lines around matches

### Integration
- Simple Go API for executing searches programmatically
- Structured result format with file paths, line numbers, and match context
- JSON output parsing for consistent data handling

## Package Structure

### Core Library

- **`lib/grep/`** - Main search functionality
  - **`request.go`** - Search request configuration and parameters
  - **`response.go`** - Structured response types and result parsing
  - **`run.go`** - Search execution and ripgrep command interface

## API Overview

### Basic Usage

```go
import "github.com/yourproject/app/lib/grep"

// Create a search request
req := &grep.Request{
    Query:     "func.*Handler",
    Directory: "./app",
    CaseSensitive: false,
}

// Execute the search
response, err := grep.Run(req)
if err != nil {
    return err
}

// Process results
for _, match := range response.Matches {
    fmt.Printf("%s:%d: %s\n", match.File, match.LineNumber, match.Line)
}
```

### Search Configuration

The module supports extensive search customization:

- **Query patterns**: Regular expressions and literal text
- **Directory targeting**: Specific paths or recursive searches
- **File filtering**: Include/exclude patterns for file types
- **Context control**: Number of lines before/after matches
- **Output options**: Various formatting and display preferences

## Requirements

### System Dependencies

- **ripgrep**: Must be installed and available in the system PATH
  - macOS: `brew install ripgrep`
  - Ubuntu/Debian: `apt install ripgrep`
  - Windows: Download from [releases page](https://github.com/BurntSushi/ripgrep/releases)

### Runtime Dependencies

- Go standard library only (no external Go dependencies)
- Works with any ripgrep version that supports JSON output

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/grep
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [ripgrep Documentation](https://github.com/BurntSushi/ripgrep/blob/master/GUIDE.md) - Complete ripgrep usage guide
- [Search Module](search.md) - Higher-level search abstraction
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
