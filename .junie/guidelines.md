# Project Guidelines

## Project Overview
Project Forge is an application that generates, manages, and grows web applications built using the Go programming language. It provides a powerful module system that allows you to add various features to your generated applications. Each module is a self-contained component that adds specific functionality to your project.

### Module System
Each module is represented by a folder in `./module` and contains a `.module.json` file that defines its properties, dependencies, and configuration. Modules can be combined freely to create the exact feature set your application needs.

### Available Modules
The following modules are available in Project Forge:

- **Android**: Enables building Android applications
- **Audit**: Provides audit logging and tracking capabilities
- **Brands**: Manages brand-specific assets and configurations
- **Core**: Provides common utilities for Go applications, including CLI actions, HTTP controllers, logging, telemetry, themes, and user management
- **Database**: Provides an API for accessing relational databases
- **DatabaseUI**: Offers a web interface for database management
- **Desktop**: Supports desktop application development
- **DocBrowse**: Enables browsing and viewing documentation
- **Export**: Handles data export functionality
- **Expression**: Provides expression parsing and evaluation
- **Filesystem**: Adds filesystem operations and management
- **Git**: Provides Git integration capabilities
- **GraphQL**: Implements GraphQL API support
- **HAR**: Handles HTTP Archive (HAR) file processing
- **Help**: Provides help and documentation system
- **iOS**: Enables building iOS applications
- **JSX**: Adds support for JSX/React development
- **Marketing**: Includes marketing-related features and tools
- **MCP**: Master Control Program for system management
- **Metamodel**: Handles data modeling and code generation
- **Migration**: Manages database schema migrations
- **MySQL**: Provides MySQL database support
- **Notarize**: Handles code signing and notarization
- **Notebook**: Implements interactive notebook functionality
- **Numeric**: Provides advanced numeric operations
- **OAuth**: Implements OAuth authentication
- **OpenAPI**: Supports OpenAPI/Swagger integration
- **Playwright**: Enables end-to-end testing using Playwright
- **Plot**: Provides data visualization and plotting
- **Postgres**: Provides PostgreSQL database support
- **Process**: Manages system processes and execution
- **Proxy**: Implements proxy server functionality
- **Queue**: Provides message queue integration
- **ReadOnlyDB**: Supports read-only database connections
- **RichEdit**: Implements rich text editing capabilities
- **Sandbox**: Provides isolated execution environment
- **Schedule**: Manages scheduled tasks and jobs
- **Scripting**: Supports custom scripting functionality
- **Search**: Implements search functionality
- **SQLite**: Provides SQLite database support
- **SQLServer**: Provides Microsoft SQL Server support
- **Task**: Manages background tasks and processing
- **ThemeCatalog**: Provides additional UI themes
- **Types**: Implements custom type system
- **Upgrade**: Manages application upgrades
- **User**: Handles user management and authentication
- **WASMClient**: Supports WebAssembly client applications
- **WASMServer**: Allows building HTTP servers that can run as WebAssembly modules
- **WebSocket**: Adds WebSocket support to your application

Each module can be enabled or disabled based on your project's requirements, and multiple modules can work together to provide comprehensive functionality.

## Project Structure
- `/app` - Core application code
- `/assets` - Static assets
- `/bin` - Build and utility scripts
- `/client` - Client-side code
- `/doc` - Project documentation
- `/test` - Test files
- `/views` - Web templates
- `/module` - Project modules
- `/.projectforge` - Project configuration files

## Project Configuration
Projects created by Project Forge include a configuration file at `.projectforge/project.json` that defines the project's settings and enabled modules. This file contains:

- **Basic Information**
  - `key`: Unique project identifier
  - `name`: Display name of the project
  - `icon`: Project icon identifier
  - `version`: Current version number
  - `package`: Go package path
  - `port`: Default HTTP port

- **Module Configuration**
  - `modules`: List of enabled Project Forge modules
  - `ignore`: Patterns for files/directories to ignore

- **Project Metadata**
  - `tags`: Project categorization tags
  - `info`: Detailed project information including:
    - Organization and author details
    - License and homepage
    - Project description
    - Bundle identifiers
    - Build-specific configurations

- **Build Settings**
  - `build`: Configuration for various build targets and platforms

Example configuration:
```json
{
  "key": "myproject",
  "name": "My Project",
  "version": "1.0.0",
  "package": "github.com/username/myproject",
  "port": 8080,
  "modules": [
    "core",
    "database",
    "websocket"
  ],
  "info": {
    "org": "myorg",
    "authorName": "Author Name",
    "license": "MIT",
    "homepage": "https://myproject.com",
    "description": "Project description"
  },
  "build": {
    "desktop": true,
    "wasm": true
  }
}
```

## Development Setup

### Prerequisites
- Go 1.24.0 or higher
- Make (for build scripts)

### Build Commands
```bash
make dev        # Start development mode with hot reload
make build      # Build debug binary
make build-release  # Build release binary
make templates  # Compile templates
make lint       # Run linter
make clean      # Clean build artifacts
```

### Testing
Tests are located in the `/test` directory. Run tests using standard Go testing commands:
```bash
go test ./...
```

## Best Practices
1. **Code Organization**
   - Follow Go standard project layout
   - Keep modules self-contained
   - Use proper package naming

2. **Development Workflow**
   - Use `make dev` for local development
   - Run linter before commits
   - Update templates when modifying views

3. **Documentation**
   - Document new features in `/doc`
   - Keep API documentation up-to-date
   - Follow existing documentation patterns

## Key Technologies
- Gorilla Mux (Web Router)
- OpenTelemetry (Observability)
- QuickTemplate (Template Engine)
- Zap (High-performance logging)
- Coral (CLI interface)
- Lo (Functional programming utilities)
- Chroma (Syntax highlighting)
- WebSocket Support
- Markdown Processing
- YAML/TOML Configuration
- Prometheus (Metrics and monitoring)
- JSON Iterator (Fast JSON processing)

## Additional Resources
- Detailed documentation in `/doc` directory
- Build scripts in `/bin` directory
- Configuration examples in project root
