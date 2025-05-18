# AGENTS.md

This file provides guidance to autonomous agents for working with code in this repository.

## Project Overview
Project Forge is an application that generates, manages, and grows web applications built using the Go programming language. It provides a powerful module system that allows you to add various features to your generated applications. Each module is a self-contained component that adds specific functionality to your project.

### Module System
Each module is represented by a folder in `./module` and contains a `.module.json` file that defines its properties, dependencies, and configuration. Modules can be combined freely to create the exact feature set your application needs.

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

## Build Commands
- Build: `make build`
- Run: `make build`, then `./build/debug/projectforge`
- Run and watch: `make dev` or `./bin/dev.sh`
- Lint: `make lint` or `./bin/check.sh`
- Format: `./bin/format.sh`
- Test all: `./bin/test.sh` or `gotestsum -- ./app/...`
- Test single: `go test ./path/to/package -run TestName`
- Test with clean cache: `./bin/test.sh -c`
- Test and watch: `./bin/test.sh -w`

## Code Style Guidelines
- Use `gofumpt` for formatting (enforced by linters)
- Import order: standard libs, third-party libs, project imports
- Max line length: 160 characters
- Error handling: check all errors, use appropriate context
- Functions: max 100 lines recommended
- Cyclomatic complexity: max 30 recommended
- Naming: PascalCase for exported, camelCase for internal
- Comments: All exported items must be documented
- Use quicktemplate for HTML templates
- Follow existing patterns when adding new code

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

## Common Development Patterns
- Controllers are in `/app/controller` and follow MVC pattern
- Templates are in `/views` and use quicktemplate syntax
- Database models typically include CRUD operations
- Error handling uses pkg/errors for stack traces
- Configuration uses environment variables with defaults
- CLI commands are defined in `/app/cmd`
- Logging uses structured logging via zap

## Testing Approach
- Unit tests use standard Go testing package
- Integration tests are in `/test` directory
- Test files are named with `_test.go` suffix
- Mock objects are used for external dependencies
- Test utilities are in `/test/util`

## Additional Resources
- Detailed documentation in `/doc` directory
- Build scripts in `/bin` directory
- Configuration examples in project root
