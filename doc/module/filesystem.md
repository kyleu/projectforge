# Filesystem

The **`filesystem`** module provides a unified abstraction layer for [Project Forge](https://projectforge.dev) applications to interact with different storage backends. It enables applications to work seamlessly with local disk storage, in-memory filesystems, and remote storage systems through a consistent API.

## Overview

This module provides:

- **Unified Interface**: Single `FileLoader` interface for all filesystem operations
- **Multiple Backends**: Support for disk-based and memory-based storage
- **File Operations**: Complete CRUD operations for files and directories
- **Advanced Features**: Recursive operations, file watching, compression, and downloads
- **Performance**: Efficient operations with streaming support for large files

## Key Features

### Storage Backends
- **Disk Storage**: Traditional filesystem operations using `afero.OsFs`
- **Memory Storage**: In-memory filesystem for testing and temporary data using `afero.MemMapFs`
- **Read-Only Mode**: Configurable read-only access for security
- **Cross-Platform**: Automatic backend selection based on platform (JS uses memory)

### File Operations
- File reading, writing, and streaming
- Directory creation and traversal
- File copying and moving
- Recursive operations with ignore patterns
- JSON file handling with type safety
- File compression and extraction

### Advanced Capabilities
- **Remote Downloads**: HTTP/HTTPS file downloads with progress tracking
- **Archive Support**: ZIP file extraction and creation
- **Tree Structures**: Generate file tree representations
- **Pattern Matching**: File filtering with glob patterns
- **Metadata Access**: File permissions, sizes, and modification times

## Package Structure

### Core Interface

- **`FileLoader`** - Primary interface for all filesystem operations
  - File I/O operations (read, write, copy, move)
  - Directory management (create, list, traverse)
  - Metadata operations (stat, permissions, exists)
  - Advanced features (download, compress, stream)

### Implementation

- **`FileSystem`** - Main implementation using [afero](https://github.com/spf13/afero)
  - Configurable storage backends (disk/memory)
  - Read-only mode support
  - Path resolution and root directory management
  - Cross-platform compatibility

### File Types and Utilities

- **`FileInfo`** - Enhanced file metadata with Project Forge-specific fields
- **`FileMode`** - Type-safe file permission handling
- **`Reader/Writer`** - Streaming interfaces for large file operations
- **`Tree`** - Hierarchical file structure representation

## Usage Examples

### Basic File Operations

```go
// Create a new filesystem rooted at ./data
fs, err := NewFileSystem("./data", false, "file")
if err != nil {
    return err
}

// Read a file
content, err := fs.ReadFile("config.json")
if err != nil {
    return err
}

// Write a file
err = fs.WriteFile("output.txt", []byte("Hello World"), 0644, false)
if err != nil {
    return err
}
```

### JSON File Handling

```go
// Write structured data as JSON
data := map[string]interface{}{
    "version": "1.0",
    "enabled": true,
}
err := fs.WriteJSONFile("config.json", data, 0644, true)
```

### Directory Operations

```go
// Create directory structure
err := fs.CreateDirectory("logs/application")

// List files with filtering
files := fs.ListFiles("./", []string{".git", "node_modules"}, logger)

// Recursive file listing with patterns
matches, err := fs.ListFilesRecursive("./src", nil, logger, "*.go")
```

### Advanced Features

```go
// Download remote file
size, err := fs.Download(ctx, "https://example.com/file.zip", "downloads/file.zip", true, logger)

// Extract ZIP archive
fileMap, err := fs.UnzipToDir("archive.zip", "extracted/")

// Copy with ignore patterns
err := fs.CopyRecursive("src/", "backup/", []string{"*.tmp", ".DS_Store"}, logger)
```

## Configuration

The filesystem can be configured through the `NewFileSystem` constructor:

- **`root`**: Base directory for all operations
- **`readonly`**: Enable read-only mode for security
- **`mode`**: Storage backend ("file", "memory", or "" for auto-detection)

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/filesystem
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [afero](https://github.com/spf13/afero) - Filesystem abstraction library providing the underlying storage implementations
- [Core Module](core.md) - Foundation module with basic utilities
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
