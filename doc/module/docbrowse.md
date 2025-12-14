# Doc Browser

The **`docbrowse`** module provides a comprehensive documentation browser for your application.
It automatically converts Markdown files into a navigable web interface with hierarchical organization and responsive navigation.

## Overview

This module transforms your `./doc` directory into a fully-featured documentation website with:

- **Automatic Navigation**: Generates hierarchical menus from filesystem structure
- **Markdown Rendering**: Converts Markdown files to clean, styled HTML
- **Content Negotiation**: Supports both HTML and JSON responses
- **Responsive Design**: Mobile-friendly documentation browsing

## Key Features

### Automatic Organization
- Scans `./doc` directory and creates structured navigation
- Extracts document titles from first `#` heading
- Supports nested folder hierarchies
- Alphabetical sorting with proper precedence

### Content Management
- Real-time Markdown to HTML conversion
- Clean, readable HTML output
- Breadcrumb navigation for deep hierarchies
- File-based content organization (no database required)

### Developer Experience
- Zero configuration required
- Supports standard Markdown syntax
- Integrates with theming system
- Content negotiation for API access

## Package Structure

### Controllers

- **`controller/docs.go`** - Main documentation serving logic
  - Route handling for `/docs/{path}` endpoints
  - Markdown file discovery and parsing
  - Menu generation from filesystem structure
  - Content negotiation (HTML/JSON responses)

### Libraries

- **`lib/doc/`** - Documentation processing utilities
  - Markdown parsing and HTML conversion
  - File tree navigation and organization
  - Title extraction from document headers
  - Menu hierarchy generation

### Views

- **`views/docs/`** - Documentation display templates
  - Clean, readable documentation layouts
  - Responsive navigation components
  - Breadcrumb trail generation
  - Mobile-optimized reading experience

## Usage

### Basic Setup

1. **Add Markdown files** to your `./doc` directory:
   ```
   doc/
   ├── getting-started.md
   ├── api/
   │   ├── overview.md
   │   └── reference.md
   └── guides/
       ├── deployment.md
       └── troubleshooting.md
   ```

2. **Access documentation** at `/docs` in your application

3. **Navigation is automatically generated** from the folder structure

### Document Structure

Each Markdown file should start with a heading:

```markdown
# Document Title

Your content here...
```

The first `#` heading becomes the document title in navigation menus.

### Organizing Content

- **Folders** create menu sections
- **Files** become individual documentation pages
- **Nested folders** create sub-menus
- **File names** determine URL structure

## API Endpoints

### GET /docs
Returns the documentation index page with full navigation menu.

### GET /docs/{path}
Serves individual documentation pages. Supports:
- **HTML**: Default browser response with full page layout
- **JSON**: Raw document data for API consumption

## Configuration

No configuration required. The module automatically:
- Discovers all `.md` files in `./doc`
- Generates navigation structure
- Serves content at documented endpoints

## Source Code

- **Repository**: https://github.com/kyleu/projectforge/tree/main/module/docbrowse
- **License**: [CC0](https://creativecommons.org/publicdomain/zero/1.0) (Public Domain)
- **Author**: Kyle U (kyle@kyleu.com)

## See Also

- [Markdown Guide](https://www.markdownguide.org) - Markdown syntax reference
- [Project Forge Documentation](https://projectforge.dev) - Complete documentation
