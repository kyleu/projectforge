# Notebook

The **`notebook`** module provides [Observable Framework](https://observablehq.com/framework) integration for your application. It enables embedded data visualization, interactive analysis, and notebook-style development workflows within your application.

## Overview

This module integrates the Observable Framework to provide:

- **Embedded Notebooks**: View and interact with Observable notebooks directly in your application
- **File Management**: Browse, edit, and manage notebook files through a web interface
- **Development Server**: Automatic Observable Framework development server on a dedicated port
- **Live Editing**: Real-time editing of notebook files with syntax highlighting
- **Process Management**: Start, stop, and monitor the Observable Framework server

## Key Features

### Observable Framework Integration
- Complete TypeScript configuration for Observable Framework
- Automatic development server startup on `[app port + 10]`
- Embedded notebook viewing at `/notebook` routes
- Support for all Observable Framework features (data loaders, components, etc.)

### File Management
- Browse notebook files through web interface
- Edit files with syntax highlighting
- Save changes with validation
- Directory navigation and file organization

### Development Experience
- Seamless integration with your application's development workflow
- Automatic server management and health checking
- Configurable base URL for flexible deployment
- Real-time status monitoring

## Package Structure

### Controllers

- **`controller/cnotebook.go`** - Main notebook controller
  - View embedded notebooks at `/notebook`
  - Browse files at `/notebook/files`
  - Edit files at `/notebook/edit`
  - Save file changes with validation

### Services

- **`lib/notebook/`** - Observable Framework service integration
  - Server process management
  - Status checking and health monitoring
  - File system operations
  - Configuration management

### Templates

- **Notebook viewing**: Embedded iframe integration
- **File browser**: Directory and file listing interface
- **File editor**: Syntax-highlighted code editor
- **Status pages**: Server status and management controls

## Configuration

### Observable Framework Setup

The module automatically configures Observable Framework with:
- TypeScript support
- Development server configuration
- Integration with your application's build system
- Proper routing for embedded viewing

## Usage

### Basic Setup

1. Add the `notebook` module to your application
2. The Observable Framework development server can be started automatically upon first use
3. Access notebooks at `/notebook` in your application
4. Manage files through `/notebook/files`

### Creating Notebooks

1. Navigate to `/notebook/files` in your application
2. Create new `.md` files for Observable notebooks
3. Use Observable Framework syntax for data visualization
4. Files are automatically served by the embedded server

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/notebook
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Observable Framework Documentation](https://observablehq.com/framework) - Complete Observable Framework guide
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
